package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/pkg/pi"
	"github.com/lapsang-boys/pippi/pkg/services"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// Default gRPC address to listen on.
	defaultStringsAddr = fmt.Sprintf("localhost:%d", services.StringsPort)
	// Default minimum string length in characters.
	defaultMinLength = 4
)

// serverCmd is the command to launch a gRPC server processing extract strings
// requests.
type serverCmd struct {
	// gRPC address to listen on.
	stringsAddr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to extract strings requests.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.stringsAddr, "addr", defaultStringsAddr, "gRPC address to listen on")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.stringsAddr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to extract
// strings.
func listen(stringsAddr string) error {
	dbg.Printf("listening on %q", stringsAddr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", stringsAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register binary parser service.
	stringspb.RegisterStringsExtractorServer(server, &stringsExtractorServer{})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// stringsExtractorServer implements stringspb.StringsExtractorServer.
type stringsExtractorServer struct{}

// ExtractStrings extracts printable strings from the given binary file.
func (s *stringsExtractorServer) ExtractStrings(ctx context.Context, req *stringspb.StringsRequest) (*stringspb.StringsReply, error) {
	if err := pi.CheckBinID(req.BinId); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("parsing ID %q", req.BinId)
	// Read contents of binary file.
	binPath, err := pi.BinPath(req.BinId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buf, err := ioutil.ReadFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Extract strings.
	// TODO: make minimum length configurable; add to strings.proto?
	infos := extractStrings(buf, defaultMinLength)
	// Send reply.
	reply := &stringspb.StringsReply{
		Strings: infos,
	}
	return reply, nil
}
