package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	"github.com/lapsang-boys/pippi/pkg/pi"
	"github.com/lapsang-boys/pippi/pkg/services"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
)

func binary(binId string) []byte {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	binPath, err := pi.BinPath(binId)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	// Read file contents.
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
	binAddr := fmt.Sprintf("localhost:%d", services.BinPort)
	file, err := Sections(binAddr, binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

func disassembly(binId string) []*disasmpb.Instruction {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	disasmObjdumpAddr := fmt.Sprintf("localhost:%d", services.DisasmObjdumpPort)
	instAddrs, err := InstAddrs(disasmObjdumpAddr, binId)
	if err != nil {
		log.Println(err)
		return nil
	}
	arch := binpb.Arch_X86_64 // TODO: make configurable.
	disasmAddr := fmt.Sprintf("localhost:%d", services.DisasmPort)
	reply, err := Disassembly(disasmAddr, binId, arch, instAddrs)
	if err != nil {
		log.Println(err)
		return nil
	}
	return reply.Insts
}

func strings(binId string) []*stringspb.StringInfo {
	if err := pi.CheckBinID(binId); err != nil {
		log.Printf("invalid binary ID %q: %v", binId, err)
		return nil
	}
	stringsAddr := fmt.Sprintf("localhost:%d", services.StringsPort)
	extractedStrings, err := Strings(stringsAddr, binId)
	if err != nil {
		log.Println(err)
		return []*stringspb.StringInfo{}
	}
	return extractedStrings
}

func listIds() []string {
	pippiCacheDir, err := pi.CacheDir()
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(pippiCacheDir)
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
