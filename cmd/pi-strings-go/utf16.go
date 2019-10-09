package main

import (
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	textunicode "golang.org/x/text/encoding/unicode"
)

// extractUTF16LittleEndianStrings extracts printable UTF-16 strings in little
// endian byte order with a minimum number of characters from the given binary.
func extractUTF16LittleEndianStrings(buf []byte, minLength int, c chan []*stringspb.StringInfo) {
	extractUTF16Strings(buf, minLength, stringspb.Encoding_UTF16LittleEndian, textunicode.LittleEndian, c)
}

// extractUTF16BigEndianStrings extracts printable UTF-16 strings in big endian
// byte order with a minimum number of characters from the given binary.
func extractUTF16BigEndianStrings(buf []byte, minLength int, c chan []*stringspb.StringInfo) {
	extractUTF16Strings(buf, minLength, stringspb.Encoding_UTF16BigEndian, textunicode.BigEndian, c)
}

// extractUTF16Strings extracts printable UTF-8 strings with a minimum number of
// characters from the given binary. The given byte order is used to decode
// uint16 values of buf.
func extractUTF16Strings(buf []byte, minLength int, encoding stringspb.Encoding, endianness textunicode.Endianness, c chan []*stringspb.StringInfo) {
	var infos []*stringspb.StringInfo
	for i := 0; i < len(buf); {
		start := uint64(i)
		s, n, ok := findUTF16String(buf[start:], minLength, endianness)
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
// at src. For the string to be valid, it must be of at least the specified
// minimum length in number of characters. The given byte order is used to
// decode uint16 values of src. The integer return value n specifies the number
// of bytes read, and the boolean return value indicates success. If an invalid
// UTF-8 encoding is encountered at the start of the given buffer, n is set to 1
// and the boolean return value is false.
func findUTF16String(src []byte, minLength int, endianness textunicode.Endianness) (s string, n uint64, ok bool) {
	// Throwaway buffer needed for encoding.Encoding.Transform.
	const maxSize = 10 * 1024 * 1024 // 10 MB
	var dst [maxSize]byte
	// TODO: add BOM check.
	//utf16BigBom := textunicode.UTF16(textunicode.BigEndian, textunicode.UseBOM)
	//utf16Big := textunicode.UTF16(textunicode.BigEndian, textunicode.IgnoreBOM)
	dec := textunicode.UTF16(endianness, textunicode.IgnoreBOM).NewDecoder()
	nDst, nSrc, _ := dec.Transform(dst[:], src, false)
	// TODO: check number of runes decoded, not number of bytes.
	if nDst > minLength {
		return string(dst[:nDst]), uint64(nSrc), true
	}
	return "", 1, false
}
