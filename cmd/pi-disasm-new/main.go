package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf"
	_ "github.com/decomp/exp/bin/pe"
	"github.com/lapsang-boys/pippi/cmd/pi-disasm-new/disasm"
	_ "github.com/lapsang-boys/pippi/cmd/pi-disasm-new/disasm/x86"
	"github.com/mewkiz/pkg/jsonutil"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "pi-disasm-new:" prefix which logs debug messages
	// to standard error.
	dbg = log.New(os.Stderr, term.CyanBold("pi-disasm-new:")+" ", 0)
	// warn is a logger with the "pi-disasm-new:" prefix which logs warning
	// messages to standard error.
	warn = log.New(os.Stderr, term.RedBold("pi-disasm-new:")+" ", 0)
)

func usage() {
	const use = `
Usage:

pi-disasm-new [OPTION]... BIN_FILE`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	binPath := flag.Arg(0)
	// Parse database.
	db, err := parseDatabase(binPath)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// Disassemble binary file.
	insts, err := disasmBinary(db, binPath)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	for _, inst := range insts {
		fmt.Println(inst)
	}
}

// A Database holds the state of a reverse engineering project.
type Database struct {
	// Addresses of instructions in the binary.
	InstAddrs []bin.Address
}

// parseDatabase parses the database of the given binary.
func parseDatabase(binPath string) (*Database, error) {
	dbPath := pathutil.TrimExt(binPath) + ".json"
	var instAddrs []bin.Address
	if err := jsonutil.ParseFile(dbPath, &instAddrs); err != nil {
		return nil, errors.WithStack(err)
	}
	db := &Database{
		InstAddrs: instAddrs,
	}
	return db, nil
}

// disasmBinary disassembles the given binary.
func disasmBinary(db *Database, binPath string) ([]disasm.Instruction, error) {
	// Parse bianry file.
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Create disassembler based on machine architecture.
	dis, err := disasm.NewDisassembler(file.Arch)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Disassemble instructions.
	var insts []disasm.Instruction
	for _, instAddr := range db.InstAddrs {
		data := file.Code(instAddr)
		inst, err := dis.DecodeInst(data)
		if err != nil {
			//fmt.Fprintln(os.Stderr, hex.Dump(data))
			warn.Printf("unable to decode instruction at address %v", instAddr)
			continue
		}
		insts = append(insts, inst)
	}
	return insts, nil
}
