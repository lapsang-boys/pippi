package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/kr/pretty"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to connect to a gRPC server to send extract strings
// requests.
type clientCmd struct {
	// gRPC address to connect to.
	Addr string
}

func (*clientCmd) Name() string {
	return "client"
}

func (*clientCmd) Synopsis() string {
	return "connect to gRPC server"
}

func (*clientCmd) Usage() string {
	const use = `
Send binary extract strings requests.

Usage:
	client [OPTION]... BIN_ID

Flags:
`
	return use[1:]
}

func (cmd *clientCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.Addr, "addr", defaultAddr, "gRPC address to connect to")
}

func (cmd *clientCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binID := f.Arg(0)

	// Connect to gRPC server.
	if err := connect(cmd.Addr, binID); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address to send an extract strings
// request.
func connect(addr, binID string) error {
	dbg.Printf("connecting to %q", addr)
	// Connect to gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()
	// Send extract strings request.
	client := stringspb.NewStringsExtractorClient(conn)
	ctx := context.Background()
	req := &stringspb.StringsRequest{
		BinId: binID,
	}
	infos, err := client.ExtractStrings(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println(infos)
	return nil
}
