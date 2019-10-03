// Package pi implements utility functions shared across Pippi.
package pi

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// CheckBinID validates the given binary ID.
func CheckBinID(binID string) error {
	if len(binID) != sha256.Size*2 {
		return errors.Errorf("invalid length of binary ID; expected %d, got %d", sha256.Size*2, len(binID))
	}
	if binID != strings.ToLower(binID) {
		return errors.Errorf("invalid binary ID; expected lowercase, got %q", binID)
	}
	const hex = "0123456789abcdef"
	for _, r := range binID {
		if !strings.ContainsRune(hex, r) {
			return errors.Errorf("invalid character in binary ID; expected hexadecimal digit, got %q", r)
		}
	}
	return nil
}

// CacheDir returns the pippi cache directory.
func CacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.WithStack(err)
	}
	pippiCacheDir := filepath.Join(cacheDir, "pippi")
	return pippiCacheDir, nil
}

// BinDir returns the project directory of the given binary ID.
func BinDir(binID string) (string, error) {
	if err := CheckBinID(binID); err != nil {
		return "", errors.WithStack(err)
	}
	pippiCacheDir, err := CacheDir()
	if err != nil {
		return "", errors.WithStack(err)
	}
	binDir := filepath.Join(pippiCacheDir, binID)
	return binDir, nil
}

// Binary executable extension.
const binExt = ".bin"

// BinPath returns the file path to the binary of the given binary ID.
func BinPath(binID string) (string, error) {
	binDir, err := BinDir(binID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	binName := binID + binExt
	binPath := filepath.Join(binDir, binName)
	return binPath, nil
}
