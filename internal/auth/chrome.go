package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"

	_ "modernc.org/sqlite"
)

// ErrNoChromeProfile is returned when Chrome's cookie DB cannot be found.
var ErrNoChromeProfile = errors.New("Chrome cookie database not found")

// ErrCookieNotFound is returned when the requested cookie is absent.
var ErrCookieNotFound = errors.New("cookie not found in Chrome")

// ChromeCookies holds the extracted Twitter auth cookies.
type ChromeCookies struct {
	AuthToken string
	CT0       string
}

// ExtractFromChrome reads auth_token and ct0 from Chrome's encrypted cookie store.
// Only works on macOS.
func ExtractFromChrome() (*ChromeCookies, error) {
	dbPath, err := chromeCookieDBPath()
	if err != nil {
		return nil, err
	}

	key, err := chromeDerivedKey()
	if err != nil {
		return nil, fmt.Errorf("derive Chrome AES key: %w", err)
	}

	// Work on a copy of the cookie DB so we don't lock Chrome's file.
	tmp, err := copyToTemp(dbPath)
	if err != nil {
		return nil, fmt.Errorf("copy cookie DB: %w", err)
	}
	defer os.Remove(tmp)

	db, err := sql.Open("sqlite", tmp)
	if err != nil {
		return nil, fmt.Errorf("open cookie DB: %w", err)
	}
	defer db.Close()

	authToken, err := queryCookie(db, key, "auth_token")
	if err != nil {
		return nil, fmt.Errorf("auth_token: %w", err)
	}
	ct0, err := queryCookie(db, key, "ct0")
	if err != nil {
		return nil, fmt.Errorf("ct0: %w", err)
	}

	return &ChromeCookies{AuthToken: authToken, CT0: ct0}, nil
}

// chromeCookieDBPath returns the path to Chrome's Default profile cookie DB.
func chromeCookieDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Try Default profile first, then Profile 1.
	candidates := []string{
		filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default", "Cookies"),
		filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Profile 1", "Cookies"),
		// Chromium
		filepath.Join(home, "Library", "Application Support", "Chromium", "Default", "Cookies"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", ErrNoChromeProfile
}

// chromeSafeStoragePassword reads the Chrome Safe Storage password from macOS Keychain.
func chromeSafeStoragePassword() (string, error) {
	out, err := exec.Command("security", "find-generic-password", "-s", "Chrome Safe Storage", "-w").Output()
	if err != nil {
		return "", fmt.Errorf("keychain read: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// chromeDerivedKey derives the 16-byte AES key from the Chrome Safe Storage password.
func chromeDerivedKey() ([]byte, error) {
	password, err := chromeSafeStoragePassword()
	if err != nil {
		return nil, err
	}
	// Chrome uses PBKDF2-SHA1 with 1003 iterations, 16-byte key, salt = "saltysalt".
	key := pbkdf2.Key([]byte(password), []byte("saltysalt"), 1003, 16, sha1.New)
	return key, nil
}

// copyToTemp copies src to a temp file and returns the temp path.
func copyToTemp(src string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	tmp, err := os.CreateTemp("", "twt-cookies-*.db")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, in); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

// queryCookie fetches and decrypts a single cookie value from the DB.
func queryCookie(db *sql.DB, key []byte, name string) (string, error) {
	// Try both host_key variants.
	hosts := []string{".x.com", ".twitter.com"}
	for _, host := range hosts {
		val, err := tryQueryCookie(db, key, name, host)
		if err == nil {
			return val, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}
	return "", ErrCookieNotFound
}

// tryQueryCookie performs the actual query for one host_key.
func tryQueryCookie(db *sql.DB, key []byte, name, hostKey string) (string, error) {
	var encrypted []byte
	err := db.QueryRow(
		`SELECT encrypted_value FROM cookies WHERE name=? AND host_key=? LIMIT 1`,
		name, hostKey,
	).Scan(&encrypted)
	if err != nil {
		return "", err
	}
	return decryptChromeValue(key, encrypted)
}

// decryptChromeValue decrypts a Chrome cookie value encrypted with AES-128-CBC.
// Chrome prepends "v10" (3 bytes) to the ciphertext.
func decryptChromeValue(key, encrypted []byte) (string, error) {
	if len(encrypted) < 3 {
		return "", errors.New("encrypted value too short")
	}
	// Strip the "v10" prefix.
	ciphertext := encrypted[3:]
	if len(ciphertext) == 0 {
		return "", errors.New("empty ciphertext after prefix strip")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext shorter than AES block size")
	}

	// IV = 16 × 0x20 (space character)
	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = 0x20
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS7 unpad.
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	// The first block (16 bytes) may contain garbage; extract the hex token.
	raw := string(plaintext)
	// Try clean extraction first (skip first block if it looks garbled).
	token := extractHexToken(raw)
	if token == "" {
		return "", errors.New("could not extract token from decrypted value")
	}
	return token, nil
}

var hexTokenRe = regexp.MustCompile(`[0-9a-f]{32,}`)

// extractHexToken extracts a lowercase hex token from the decrypted string.
func extractHexToken(s string) string {
	// Also try the full string first (for non-hex cookies like ct0).
	// ct0 is alphanumeric, auth_token is pure hex.
	s = strings.TrimSpace(s)

	// Strip null bytes and control characters.
	var cleaned strings.Builder
	for _, r := range s {
		if r >= 0x20 && r < 0x7f {
			cleaned.WriteRune(r)
		}
	}
	result := strings.TrimSpace(cleaned.String())

	// If the value looks like a raw token (no padding artifacts), return it.
	if len(result) >= 32 {
		return result
	}

	// Fall back to regex hex extraction.
	match := hexTokenRe.FindString(s)
	return match
}

// pkcs7Unpad removes PKCS7 padding from a byte slice.
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > aes.BlockSize {
		return nil, fmt.Errorf("invalid PKCS7 pad length: %d", padLen)
	}
	if padLen > len(data) {
		return nil, errors.New("pad length exceeds data length")
	}
	// Verify all pad bytes.
	for i := len(data) - padLen; i < len(data); i++ {
		if data[i] != byte(padLen) {
			return nil, errors.New("invalid PKCS7 padding")
		}
	}
	return data[:len(data)-padLen], nil
}
