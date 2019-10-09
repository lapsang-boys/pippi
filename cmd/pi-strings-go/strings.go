package main

import (
	"unicode"
	"unicode/utf8"

	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	textunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// extractStrings extracts printable strings with a minimum number of characters
// from the given binary.
func extractStrings(buf []byte, minLength int) []*stringspb.StringInfo {
	c := make(chan []*stringspb.StringInfo)
	encs := []struct {
		enc stringspb.Encoding
		dec *encoding.Decoder
	}{
		// UTF-8
		{
			enc: stringspb.Encoding_UTF8,
			dec: textunicode.UTF8.NewDecoder(),
		},
		// UTF-16
		// * UTF-16 (big endian)
		{
			enc: stringspb.Encoding_UTF16BigEndian,
			dec: textunicode.UTF16(textunicode.BigEndian, textunicode.IgnoreBOM).NewDecoder(),
		},
		// * UTF-16 (big endian, with BOM)
		{
			enc: stringspb.Encoding_UTF16BigEndianBOM,
			dec: textunicode.UTF16(textunicode.BigEndian, textunicode.ExpectBOM).NewDecoder(),
		},
		// * UTF-16 (little endian)
		{
			enc: stringspb.Encoding_UTF16LittleEndian,
			dec: textunicode.UTF16(textunicode.LittleEndian, textunicode.IgnoreBOM).NewDecoder(),
		},
		// * UTF-16 (little endian, with BOM)
		{
			enc: stringspb.Encoding_UTF16LittleEndianBOM,
			dec: textunicode.UTF16(textunicode.LittleEndian, textunicode.ExpectBOM).NewDecoder(),
		},
		// UTF-32
		// * UTF-32 (big endian)
		{
			enc: stringspb.Encoding_UTF32BigEndian,
			dec: utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewDecoder(),
		},
		// * UTF-32 (big endian, with BOM)
		{
			enc: stringspb.Encoding_UTF32BigEndianBOM,
			dec: utf32.UTF32(utf32.BigEndian, utf32.ExpectBOM).NewDecoder(),
		},
		// * UTF-32 (little endian)
		{
			enc: stringspb.Encoding_UTF32LittleEndian,
			dec: utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewDecoder(),
		},
		// * UTF-32 (little endian, with BOM)
		{
			enc: stringspb.Encoding_UTF32LittleEndianBOM,
			dec: utf32.UTF32(utf32.LittleEndian, utf32.ExpectBOM).NewDecoder(),
		},
		// * Big-5
		{
			enc: stringspb.Encoding_Big5,
			dec: traditionalchinese.Big5.NewDecoder(),
		},
		// TODO: consider adding support for ISO-88xx-xx encodings (e.g. 8859-1
		// for latin)?
		//
		//    * https://godoc.org/golang.org/x/text/encoding/charmap
		//
		// TODO: add support for more encodings? got to catch 'em all!
		//
		//    * https://godoc.org/golang.org/x/text/encoding/japanese
		//    * https://godoc.org/golang.org/x/text/encoding/korean
		//    * https://godoc.org/golang.org/x/text/encoding/simplifiedchinese
	}
	for _, enc := range encs {
		go extractEncStrings(buf, minLength, enc.enc, enc.dec, c)
	}
	var infos []*stringspb.StringInfo
	for range encs {
		infos = append(infos, <-c...)
	}
	return infos
}

// extractEncStrings extracts printable strings of the given encoding with a
// minimum number of characters from the given binary.
func extractEncStrings(buf []byte, minLength int, encoding stringspb.Encoding, dec *encoding.Decoder, c chan []*stringspb.StringInfo) {
	var infos []*stringspb.StringInfo
	for i := 0; i < len(buf); {
		start := uint64(i)
		s, n, ok := findEncString(buf[start:], minLength, dec)
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

// findEncString tries to locate the longest printable string starting at src,
// decoding from dec. For the string to be valid, it must be of at least the
// specified minimum length in number of characters. The integer return value n
// specifies the number of bytes read, and the boolean return value indicates
// success. If an invalid encoding is encountered at the start of the given
// buffer, n is set to 1 and the boolean return value is false.
func findEncString(src []byte, minLength int, dec *encoding.Decoder) (s string, n uint64, ok bool) {
	// Throwaway buffer needed for encoding.Encoding.Transform.
	const maxSize = 10 * 1024 * 1024 // 10 MB
	var dst [maxSize]byte
	nDst, nSrc, _ := dec.Transform(dst[:], src, false)
	if nDst > minLength {
		// Check number of runes decoded, not just number of bytes.
		d := dst[:nDst]
		if utf8.Valid(d) && utf8.RuneCount(d) > minLength {
			s := string(d)
			if valid(s) && prefixGraphicCount(s) > uint64(minLength) {
				return s, uint64(nSrc), true
			}
		}
	}
	return "", 1, false
}

// valid reports whether the given string is valid UTF-8 without any Unicode
// replacement characters.
func valid(s string) bool {
	for _, r := range s {
		if r == utf8.RuneError {
			return false
		}
	}
	return true
}

// prefixGraphicCount returns the number of graphic Unicode code points at the start
// of the given string.
func prefixGraphicCount(s string) uint64 {
	n := uint64(0)
	for _, r := range s {
		if !unicode.IsGraphic(r) {
			break
		}
		n++
	}
	return n
}
