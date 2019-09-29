package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	"github.com/lapsang-boys/pippi/pkg/pi"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
)

func binary(binId string) []byte {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
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
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	file, err := Sections("localhost:1200", binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

func disassembly(binId string) *disasmpb.DisassembleReply {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	reply, err := Disassembly("localhost:1300", binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return reply
}

func strings(binId string) []*stringspb.StringInfo {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	extractedStrings, err := Strings("localhost:1400", binId)
	if err != nil {
		log.Println(err)
		return []*stringspb.StringInfo{}
	}
	return extractedStrings
}

func listIds() []string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return []string{}
	}
	cacheDir = filepath.Join(cacheDir, "pippi")
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		log.Fatal(err)
	}

	var ret = make([]string, 0, len(files))
	for _, f := range files {
		ret = append(ret, f.Name())
	}
	return ret
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
		Colour: "#ffffff",
	})
	app.Bind(binary)
	app.Bind(disassembly)
	app.Bind(sections)
	app.Bind(strings)
	app.Bind(listIds)
	app.Run()
}
