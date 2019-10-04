//+build ignore

//go:generate go run gen.go

package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"unicode/utf16"
)

const (
	// ref: https://en.wikipedia.org/wiki/The_quick_brown_fox_jumps_over_the_lazy_dog
	refEN = "The quick brown fox jumps over the lazy dog"
	// ref: https://ja.wikipedia.org/wiki/The_quick_brown_fox_jumps_over_the_lazy_dog
	//
	// Quick brown eagle :)
	refJP = "素早い茶色の狐はのろまな犬どもを飛び越えた"
	// ref: https://ru.wikipedia.org/wiki/The_quick_brown_fox_jumps_over_the_lazy_dog
	//
	// Nible brown fox.
	refRU = "Шустрая бурая лисица прыгает через ленивого пса"
	// ref: https://zh.wikipedia.org/wiki/The_quick_brown_fox_jumps_over_the_lazy_dog
	//
	// Fast fox.
	refZH = "快狐跨懒狗"
)

func main() {
	// ASCII encoding.
	genASCII("ascii_en.bin", refEN)
	// UTF-8 encoding.
	genUTF8("utf8_en.bin", refEN)
	genUTF8("utf8_jp.bin", refJP)
	genUTF8("utf8_ru.bin", refRU)
	genUTF8("utf8_zh.bin", refZH)
	// UTF-16 encoding.
	// * Little endian.
	genUTF16("utf16_little_endian_en.bin", refEN, binary.LittleEndian)
	genUTF16("utf16_little_endian_jp.bin", refJP, binary.LittleEndian)
	genUTF16("utf16_little_endian_ru.bin", refRU, binary.LittleEndian)
	genUTF16("utf16_little_endian_zh.bin", refZH, binary.LittleEndian)
	// * Big endian.
	genUTF16("utf16_big_endian_en.bin", refEN, binary.BigEndian)
	genUTF16("utf16_big_endian_jp.bin", refJP, binary.BigEndian)
	genUTF16("utf16_big_endian_ru.bin", refRU, binary.BigEndian)
	genUTF16("utf16_big_endian_zh.bin", refZH, binary.BigEndian)
	// Null-terminated encoding.
	genNullTerm("null_terminated_en.bin", refEN) // Null-terminated ASCII
	genNullTerm("null_terminated_jp.bin", refJP) // Null-terminated UTF-8
	genNullTerm("null_terminated_ru.bin", refRU) // Null-terminated UTF-8
	genNullTerm("null_terminated_zh.bin", refZH) // Null-terminated UTF-8
}

// genASCII generates an ASCII encoded string based on the given reference
// string. The file contents is stored to the given path.
func genASCII(path, s string) {
	mustWriteFile(path, []byte(s))
}

// genUTF8 generates a UTF-8 encoded string based on the given reference string.
// The file contents is stored to the given path.
func genUTF8(path, s string) {
	mustWriteFile(path, []byte(s))
}

// genUTF16 generates a UTF-16 encoded string with the specified byte ordering
// based on the given reference string. The file contents is stored to the given
// path.
func genUTF16(path, s string, order binary.ByteOrder) {
	bs := utf16.Encode([]rune(s))
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, order, bs); err != nil {
		panic(err)
	}
	mustWriteFile(path, buf.Bytes())
}

// genNullTerm generates a NULL-terminated string based on the given reference
// string. The file contents is stored to the given path.
func genNullTerm(path, s string) {
	mustWriteFile(path, []byte(s+"\x00"))
}

// mustWriteFile writes data to the file of the given path. It panics on
// failure.
func mustWriteFile(path string, data []byte) {
	log.Printf("creating %q", path)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		panic(err)
	}
}
