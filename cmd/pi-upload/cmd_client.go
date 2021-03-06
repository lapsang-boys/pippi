package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/pkg/pi"
	uploadpb "github.com/lapsang-boys/pippi/proto/upload"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to upload a binary file.
type clientCmd struct {
	// gRPC address to connect to.
	uploadAddr string
}

func (*clientCmd) Name() string {
	return "client"
}

func (*clientCmd) Synopsis() string {
	return "connect to gRPC server"
}

func (*clientCmd) Usage() string {
	const use = `
Send binary file parse request.

Usage:
	client [OPTION]... FILE

Flags:
`
	return use[1:]
}

func (cmd *clientCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.uploadAddr, "addr", defaultUploadAddr, "gRPC address to connect to")
}

func (cmd *clientCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binPath := f.Arg(0)

	// Connect to gRPC server.
	if err := connect(cmd.uploadAddr, binPath); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func newRequest(binPath string) (*uploadpb.UploadRequest, error) {
	buf, err := ioutil.ReadFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req := &uploadpb.UploadRequest{
		Filename: binPath,
		Hash:     pi.BinID(buf),
		Content:  buf,
	}
	return req, nil
}

// connect connects to the given gRPC address for incoming requests to parse
// binary files.
func connect(uploadAddr, binPath string) error {
	dbg.Printf("connecting to %q", uploadAddr)
	// Launch gRPC server.
	conn, err := grpc.Dial(uploadAddr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := uploadpb.NewUploadClient(conn)
	ctx := context.Background()

	req, err := newRequest(binPath)
	if err != nil {
		return errors.WithStack(err)
	}

	reply, err := client.Upload(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(reply.Id)

	return nil
}
