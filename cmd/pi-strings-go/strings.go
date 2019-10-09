package main

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	textunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

// extractStrings extracts printable strings with a minimum number of characters
// from the given binary.
func extractStrings(buf []byte, minLength int) []*stringspb.StringInfo {
	c := make(chan []*stringspb.StringInfo)
	encs := []struct {
		enc stringspb.Encoding
		dec *encoding.Decoder
	}{
		// * UTF-8
		{
			enc: stringspb.Encoding_UTF8,
			dec: textunicode.UTF8.NewDecoder(),
		},
		// * UTF-16
		//    - big endian
		{
			enc: stringspb.Encoding_UTF16BigEndian,
			dec: textunicode.UTF16(textunicode.BigEndian, textunicode.IgnoreBOM).NewDecoder(),
		},
		//    - big endian, with BOM
		{
			enc: stringspb.Encoding_UTF16BigEndianBOM,
			dec: textunicode.UTF16(textunicode.BigEndian, textunicode.ExpectBOM).NewDecoder(),
		},
		//    - little endian
		{
			enc: stringspb.Encoding_UTF16LittleEndian,
			dec: textunicode.UTF16(textunicode.LittleEndian, textunicode.IgnoreBOM).NewDecoder(),
		},
		//    - little endian, with BOM
		{
			enc: stringspb.Encoding_UTF16LittleEndianBOM,
			dec: textunicode.UTF16(textunicode.LittleEndian, textunicode.ExpectBOM).NewDecoder(),
		},
		// * UTF-32
		//    - big endian
		{
			enc: stringspb.Encoding_UTF32BigEndian,
			dec: utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewDecoder(),
		},
		//    - big endian, with BOM
		{
			enc: stringspb.Encoding_UTF32BigEndianBOM,
			dec: utf32.UTF32(utf32.BigEndian, utf32.ExpectBOM).NewDecoder(),
		},
		//    - little endian
		{
			enc: stringspb.Encoding_UTF32LittleEndian,
			dec: utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewDecoder(),
		},
		//    - little endian, with BOM
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

		// TODO: add support for more encodings? got to catch 'em all!
		//
		//    * https://godoc.org/golang.org/x/text/encoding/japanese
		//    * https://godoc.org/golang.org/x/text/encoding/korean
		//    * https://godoc.org/golang.org/x/text/encoding/simplifiedchinese

		// TODO: add NULL-terminated string encoding.

		// TODO: add length-prefixed string encoding.
	}
	for _, enc := range encs {
		go extractEncStrings(buf, minLength, enc.enc, enc.dec, c)
	}
	var infos []*stringspb.StringInfo
	for range encs {
		infos = append(infos, <-c...)
	}
	// Sort results.
	sort.Slice(infos, func(i, j int) bool {
		if infos[i].Location < infos[j].Location {
			return true
		}
		return infos[i].Encoding < infos[j].Encoding
	})
	return infos
}

// extractEncStrings extracts printable strings of the given encoding with a
// minimum number of characters from the given binary. Results are sent on the
// channel c.
func extractEncStrings(buf []byte, minLength int, encoding stringspb.Encoding, dec *encoding.Decoder, c chan []*stringspb.StringInfo) {
	var infos []*stringspb.StringInfo
	for i := 0; i < len(buf); {
		s, n, ok := findEncString(buf[i:], minLength, dec)
		if !ok {
			i++
			continue
		}
		start := uint64(i)
		i += n
		info := &stringspb.StringInfo{
			Location:  start,
			RawString: s,
			Size:      uint64(n),
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
// buffer or the located string is too short, n is set to 1 and the boolean
// return value is false.
func findEncString(src []byte, minLength int, dec *encoding.Decoder) (s string, n int, ok bool) {
	dst := &strings.Builder{}
	nchars := 0
	for n = 0; n < len(src); {
		r, nSrc := decodeRune(src[n:], dec)
		if r == unicode.ReplacementChar {
			break
		}
		if !unicode.IsGraphic(r) {
			break
		}
		n += nSrc
		dst.WriteRune(r)
		nchars++
	}
	// Check length.
	if nchars < minLength {
		return "", 1, false
	}
	return dst.String(), n, true
}

// decodeRune tries to decode a single rune from src using the given text
// transformer. The rune may consist of several surrogate runes. The returned
// rune is either a valid rune or unicode.ReplacementChar, and the returned
// integer indicates the number of source bytes read to decode the rune.
func decodeRune(src []byte, t transform.Transformer) (rune, int) {
	const (
		// TODO: verify that 4 bytes is enough to store any Unicode code point in
		// UTF-8.
		maxDstSize = 4
		// TODO: check if any encoding requires more than 4 bytes to encode a
		// single rune (including surrogate pairs).
		maxSrcSize = 4
	)
	var dst [maxDstSize]byte
	for n := maxSrcSize; n >= 1; n-- {
		nDst, nSrc, err := t.Transform(dst[:], src[:n], false)
		if err != nil {
			continue
		}
		if n != nSrc {
			continue
		}
		d := dst[:nDst]
		if utf8.RuneCount(d) == 1 {
			r, _ := utf8.DecodeRune(d)
			return r, n
		}
	}
	return unicode.ReplacementChar, 1
}
