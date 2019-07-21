//go:generate protoc -I ../../alvriket --go_out=plugins=grpc:../../alvriket ../../alvriket/alvriket.proto

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf" // register ELF decoder
	"github.com/kr/pretty"
	"github.com/lapsang-boys/pippi/alvriket"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// dbg is a logger with the "alvriket:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.CyanBold("alvriket:")+" ", 0)
	// warn is a logger with the "alvriket:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("alvriket:")+" ", 0)
)

func usage() {
	const use = `
Parse ELF files.

Usage:

	alvriket [OPTION]... FILE

Flags:
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	var (
		// Address to listen on for gRPC.
		grpcAddr string
	)
	flag.Usage = usage
	flag.StringVar(&grpcAddr, "addr", "", `address to listen on for gRPC (e.g. ":50051")`)
	flag.Parse()

	switch {
	case len(grpcAddr) != 0:
		// Launch grpc server.
		if err := listen(grpcAddr); err != nil {
			log.Fatalf("%+v", err)
		}
	default:
		// Path to binary executable.
		var binPath string
		switch flag.NArg() {
		case 0:
			// Read from standard input.
			binPath = "-"
		case 1:
			binPath = flag.Arg(0)
		default:
			flag.Usage()
			os.Exit(1)
		}

		// Parse ELF file.
		file, err := bin.ParseFile(binPath)
		if err != nil {
			log.Fatalf("unable to parse ELF file %q: %+v", binPath, err)
		}

		// Pretty-print to standard output.
		pretty.Println(file)
	}
}

// listen listens on the given gRPC address for incoming requests to parse ELF
// files.
func listen(grpcAddr string) error {
	dbg.Printf("listening on %q", grpcAddr)
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	s := grpc.NewServer()
	alvriket.RegisterELFParserServer(s, &elfParserServer{})
	if err := s.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// elfParserServer implements alvriket.ELFParserServer.
type elfParserServer struct{}

// ParseELF parses the given ELF file.
func (elfParserServer) ParseELF(ctx context.Context, req *alvriket.ParseELFRequest) (*alvriket.ParseELFReply, error) {
	dbg.Printf("parsing %q", req.ElfPath)
	// Parse ELF file.
	file, err := bin.ParseFile(req.ElfPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &alvriket.ParseELFReply{
		Nsects: int32(len(file.Sections)),
	}
	return reply, nil
}
