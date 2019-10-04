package main

import (
	"io/ioutil"
	"testing"

	stringspb "github.com/lapsang-boys/pippi/proto/strings"
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

func TestExtractStrings(t *testing.T) {
	golden := []struct {
		path string
		want *stringspb.StringInfo
	}{
		{
			path: "../../testdata/strings/ascii_en.bin",
			want: &stringspb.StringInfo{
				Location:  0,
				RawString: refEN,
				Size:      uint64(len(refEN)),
				Encoding:  stringspb.Encoding_UTF8, // TODO: update encoding to ASCII (if only containing printable ASCII characters)
				//Encoding:  stringspb.Encoding_ASCII,
			},
		},
	}
	for _, g := range golden {
		// TODO: test different min lengths?
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Errorf("unable to read file %q; %v", g.path, err)
			continue
		}
		infos := extractStrings(buf, defaultMinLength)
		if !contains(infos, g.want) {
			t.Errorf("extracted string missing; expected string %q at location %d with encoding %v and size %d not extracted", g.want.RawString, g.want.Location, g.want.Encoding, g.want.Size)
		}
	}
}

func contains(infos []*stringspb.StringInfo, want *stringspb.StringInfo) bool {
	for _, info := range infos {
		if equal(info, want) {
			return true
		}
	}
	return false
}

func equal(a, b *stringspb.StringInfo) bool {
	if a.Location != b.Location {
		return false
	}
	if a.RawString != b.RawString {
		return false
	}
	if a.Size != b.Size {
		return false
	}
	if a.Encoding != b.Encoding {
		return false
	}
	return true
}
