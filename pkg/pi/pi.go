// Package pi implements utility functions shared across Pippi.
package pi

import (
	"crypto/sha256"
	"strings"

	"github.com/pkg/errors"
)

// CheckBinID validates the given binary ID.
func CheckBinID(id string) error {
	if len(id) != sha256.Size*2 {
		return errors.Errorf("invalid length of binary ID; expected %d, got %d", sha256.Size*2, len(id))
	}
	if id != strings.ToLower(id) {
		return errors.Errorf("invalid binary ID; expected lowercase, got %q", id)
	}
	const hex = "0123456789abcdef"
	for _, r := range id {
		if !strings.ContainsRune(hex, r) {
			return errors.Errorf("invalid rune in binary ID; expected hexadecimal digit, got %q", r)
		}
	}
	return nil
}
