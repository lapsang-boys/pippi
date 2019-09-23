package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
)

func binary(binId string) []byte {
	const (
		ext = ".bin"
	)

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	cacheDir = filepath.Join(cacheDir, "pippi")
	// Read file contents.
	binName := binId + ext
	binPath := filepath.Join(cacheDir, binId, binName)
	binData, err := ioutil.ReadFile(binPath)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return binData
}

func sections(binId string) *binpb.File {
	file, err := Sections("localhost:1200", binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

func disassembly(binId string) *disasmpb.DisassembleReply {
	reply, err := Disassembly("localhost:1300", binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return reply
}

func strings(binId string) []*stringspb.StringInfo {
	extractedStrings, err := Strings("localhost:1400", binId)
	if err != nil {
		log.Println(err)
		return []*stringspb.StringInfo{}
	}
	return extractedStrings
}

func main() {
	js := mewn.String("./frontend/dist/my-app/main-es2015.js")
	css := mewn.String("./frontend/dist/my-app/styles.css")

	go recvUploads()
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "Pippi",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(binary)
	app.Bind(disassembly)
	app.Bind(sections)
	app.Bind(strings)
	app.Run()
}
