package main

import (
	"encoding/binary"
	"unicode"
	"unicode/utf8"

	stringspb "github.com/lapsang-boys/pippi/proto/strings"
)

// extractStrings extracts printable strings with a minimum number of characters
// from the given binary.
func extractStrings(buf []byte, minLength int) []*stringspb.StringInfo {
	c := make(chan []*stringspb.StringInfo)
	fs := []func(buf []byte, minLength int, c chan []*stringspb.StringInfo){
		extractUTF8Strings,
		extractUTF16LittleEndianStrings,
		extractUTF16BigEndianStrings,
	}
	for _, f := range fs {
		go f(buf, minLength, c)
	}
	var infos []*stringspb.StringInfo
	for range fs {
		infos = append(infos, <-c...)
	}
	return infos
}

// extractUTF8Strings extracts printable UTF-8 strings with a minimum number of
// characters from the given binary.
func extractUTF8Strings(buf []byte, minLength int, c chan []*stringspb.StringInfo) {
	var infos []*stringspb.StringInfo
	for i := 0; i < len(buf); {
		start := uint64(i)
		s, n, ok := findUTF8String(buf[start:], minLength)
		i += int(n)
		if !ok {
			continue
		}
		info := &stringspb.StringInfo{
			Location:  start,
			RawString: s,
			Size:      n,
			Encoding:  stringspb.Encoding_UTF8,
		}
		infos = append(infos, info)
	}
	c <- infos
}

// findUTF8String tries to locate the longest printable UTF-8 string starting at
// buf. For the string to be valid, it must be of at least the specified minimum
// length in number of characters. The integer return value n specifies the
// number of bytes read, and the boolean return value indicates success. If an
// invalid UTF-8 encoding is encountered at the start of the given buffer, n is
// set to 1 and the boolean return value is false.
func findUTF8String(buf []byte, minLength int) (s string, n uint64, ok bool) {
	nchars := 0
	for i := 0; i < len(buf); {
		r, size := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError {
			// invalid UTF-8.
			break
		}
		if !unicode.IsGraphic(r) {
			// non-printable char.
			break
		}
		// found valid char.
		i += size
		n = uint64(i)
		nchars++
	}
	ok = nchars > minLength
	if !ok {
		return "", 1, false
	}
	return string(buf[:n]), n, true
}

// extractUTF16LittleEndianStrings extracts printable UTF-16 strings in little
// endian byte order with a minimum number of characters from the given binary.
func extractUTF16LittleEndianStrings(buf []byte, minLength int, c chan []*stringspb.StringInfo) {
	extractUTF16Strings(buf, minLength, stringspb.Encoding_UTF16LittleEndian, binary.LittleEndian, c)
}

// extractUTF16BigEndianStrings extracts printable UTF-16 strings in big endian
// byte order with a minimum number of characters from the given binary.
func extractUTF16BigEndianStrings(buf []byte, minLength int, c chan []*stringspb.StringInfo) {
	extractUTF16Strings(buf, minLength, stringspb.Encoding_UTF16BigEndian, binary.BigEndian, c)
}

// extractUTF16Strings extracts printable UTF-8 strings with a minimum number of
// characters from the given binary. The given byte order is used to decode
// uint16 values of buf.
func extractUTF16Strings(buf []byte, minLength int, encoding stringspb.Encoding, order binary.ByteOrder, c chan []*stringspb.StringInfo) {
	var infos []*stringspb.StringInfo
	for i := 0; i < len(buf); {
		start := uint64(i)
		s, n, ok := findUTF16String(buf[start:], minLength, order)
		i += int(n)
		if !ok {
			continue
		}
		info := &stringspb.StringInfo{
			Location:  start,
			RawString: s,
			Size:      n,
			Encoding:  encoding,
		}
		infos = append(infos, info)
	}
	c <- infos
}

// findUTF16String tries to locate the longest printable UTF-8 string starting
// at buf. For the string to be valid, it must be of at least the specified
// minimum length in number of characters. The given byte order is used to
// decode uint16 values of buf. The integer return value n specifies the number
// of bytes read, and the boolean return value indicates success. If an invalid
// UTF-8 encoding is encountered at the start of the given buffer, n is set to 1
// and the boolean return value is false.
func findUTF16String(buf []byte, minLength int, order binary.ByteOrder) (s string, n uint64, ok bool) {
	nchars := 0
	for i := 0; i < len(buf); {
		r, size := utf16DecodeRuneWithOrder(buf[i:], order)
		if r == unicode.ReplacementChar {
			// invalid UTF-16.
			break
		}
		if !unicode.IsGraphic(r) {
			// non-printable char.
			break
		}
		// found valid char.
		i += size
		n = uint64(i)
		nchars++
		s += string(r) // TODO: check if we are allowed to convert runes to strings like this.
	}
	ok = nchars > minLength
	if !ok {
		return "", 1, false
	}
	return s, n, true
}
