package main

import (
	"encoding/binary"
	"unicode"
	"unicode/utf16"
)

// utf16DecodeRuneWithOrder unpacks the first UTF-16 encoding in p and returns
// the rune and its width in bytes. The given byte order is used to decode
// uint16 values of p. If p is empty it returns (ReplacementChar, 0). Otherwise,
// if the encoding is invalid, it returns (ReplacementChar, 1)
func utf16DecodeRuneWithOrder(p []byte, order binary.ByteOrder) (rune, int) {
	var bs [2]uint16
	n := 0
	if len(p) >= 2 {
		bs[0] = order.Uint16(p)
		n++
	}
	if len(p) >= 4 {
		bs[1] = order.Uint16(p[2:])
		n++
	}
	return utf16DecodeRune(bs[:n])
}

// Copied from unicode/utf16 package of the Go 1.13 stdlib.
const (
	// 0xd800-0xdc00 encodes the high 10 bits of a pair.
	// 0xdc00-0xe000 encodes the low 10 bits of a pair.
	// the value is those 20 bits plus 0x10000.
	surr1 = 0xd800
	surr2 = 0xdc00
	surr3 = 0xe000

	surrSelf = 0x10000
)

// utf16DecodeRune unpacks the first UTF-16 encoding in p and returns the rune
// and its width in bytes. If p is empty it returns (ReplacementChar, 0).
// Otherwise, if the encoding is invalid, it returns (ReplacementChar, 1)
//
// Copied with minor modifications from unicode/utf16.Decode function of the Go
// 1.13 stdlib.
//
// Changes made:
//    * Unpack only first rune in p (which may consist of two surrogate runes)
//    * Return the number of bytes read to decode the rune.
func utf16DecodeRune(p []uint16) (rune, int) {
	if len(p) < 1 {
		return unicode.ReplacementChar, 0
	}
	switch r := p[0]; {
	case r < surr1, surr3 <= r:
		// normal rune
		return rune(r), 2
	case surr1 <= r && r < surr2 && len(p) > 1 && surr2 <= p[1] && p[1] < surr3:
		// valid surrogate sequence
		return utf16.DecodeRune(rune(r), rune(p[1])), 4
	default:
		// invalid surrogate sequence
		return unicode.ReplacementChar, 1
	}
}
