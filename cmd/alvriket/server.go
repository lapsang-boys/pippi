package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf" // register ELF decoder
	_ "github.com/decomp/exp/bin/pe"  // register PE decoder
	"github.com/google/subcommands"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default gRPC address to listen on.
	grpcAddr = ":1234"
	// Extension of binary files.
	ext = ".bin"
)

// serverCmd is the command to launch a gRPC server processing parse binary
// requests.
type serverCmd struct {
	// gRPC address to listen on.
	Addr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to binary file parse request.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.Addr, "addr", grpcAddr, "gRPC address to listen on")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.Addr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to parse binary
// files.
func listen(addr string) error {
	dbg.Printf("listening on %q", addr)
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.WithStack(err)
	}
	cacheDir = filepath.Join(cacheDir, "pippi")
	// Launch gRPC server.
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register binary parser service.
	binpb.RegisterBinaryParserServer(server, &binParserServer{cacheDir: cacheDir})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// binParserServer implements binpb.BinaryParserServer.
type binParserServer struct {
	// Cache directory of pippi.
	cacheDir string
}

// ParseBinary parses the given binary file.
func (s *binParserServer) ParseBinary(ctx context.Context, req *binpb.ParseBinaryRequest) (*binpb.ParseBinaryReply, error) {
	if err := validateID(req.BinId); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("parsing ID %q", req.BinId)
	// Parse binary file.
	binName := filepath.Base(req.BinId) + ext
	binPath := filepath.Join(s.cacheDir, binName)
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &binpb.ParseBinaryReply{}
	for _, sect := range file.Sections {
		section := &binpb.Section{
			Name:     sect.Name,
			Addr:     uint64(sect.Addr),
			Offset:   sect.Offset,
			Length:   uint64(len(sect.Data)),
			FileSize: uint64(sect.FileSize),
			MemSize:  uint64(sect.MemSize),
		}
		if sect.Perm&bin.PermR != 0 {
			section.Perms = append(section.Perms, binpb.Perm_PermR)
		}
		if sect.Perm&bin.PermW != 0 {
			section.Perms = append(section.Perms, binpb.Perm_PermW)
		}
		if sect.Perm&bin.PermX != 0 {
			section.Perms = append(section.Perms, binpb.Perm_PermX)
		}
		reply.Sections = append(reply.Sections, section)
	}
	return reply, nil
}

// validateID validates the given binary ID.
func validateID(id string) error {
	if sha256.Size*2 != len(id) {
		return errors.Errorf("invalid length of binary ID; expected %d, got %d", sha256.Size*2, len(id))
	}
	s := strings.ToLower(id)
	if s != id {
		return errors.Errorf("invalid binary ID; expected lowercase, got %q", id)
	}
	const hex = "0123456789abcdef"
	for _, r := range id {
		if !strings.ContainsRune(hex, r) {
			return errors.Errorf("invalid rune in binary ID; expected hexadecimal digit, got %q", r)
		}
	}
	return nil
}
