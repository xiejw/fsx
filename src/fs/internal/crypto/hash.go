package crypto

import (
	"crypto/sha256"
	"io"
	"os"
)

// Returns the Sha256Sum of the file given by `path`.
func Sha256Sum(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return sha256Sum(f)
}

// Internal helper method for testing.
func sha256Sum(reader io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
