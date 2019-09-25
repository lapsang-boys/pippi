// Package binpbx provides bin.proto utility functions.
package binpbx

import (
	"context"
	"log"
	"os"

	binpb "github.com/lapsang-boys/pippi/proto/bin"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// dbg is a logger with the "binpbx:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.CyanBold("binpbx:")+" ", 0)
	// warn is a logger with the "binpbx:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("binpbx:")+" ", 0)
)

// ParseFile parses the given binary file on the specified bin gRPC server.
func ParseFile(binAddr, binID string) (*binpb.File, error) {
	dbg.Printf("connecting to %q", binAddr)
	// Connect to the bin gRPC server.
	conn, err := grpc.Dial(binAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()
	// Send binary parsing request.
	client := binpb.NewBinaryParserClient(conn)
	ctx := context.Background()
	req := &binpb.ParseBinaryRequest{
		BinId: binID,
	}
	file, err := client.ParseBinary(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return file, nil
}
