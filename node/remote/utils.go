package remote

import (
	"context"
	"crypto/tls"
	"regexp"
	"strconv"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/credentials"

	"github.com/cosmos/cosmos-sdk/codec"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	HTTPProtocols = regexp.MustCompile("https?://")
)

// GetHeightRequestContext adds the height to the context for querying the state at a given height
func GetHeightRequestContext(context context.Context, height int64) context.Context {
	return metadata.AppendToOutgoingContext(
		context,
		grpctypes.GRPCBlockHeightHeader,
		strconv.FormatInt(height, 10),
	)
}

// MustCreateGrpcConnection creates a new gRPC connection using the provided configuration and panics on error
func MustCreateGrpcConnection(cfg *GRPCConfig, cdc codec.Codec) *grpc.ClientConn {
	grpConnection, err := CreateGrpcConnection(cfg, cdc)
	if err != nil {
		panic(err)
	}
	return grpConnection
}

// CreateGrpcConnection creates a new gRPC client connection from the given configuration
func CreateGrpcConnection(cfg *GRPCConfig, cdc codec.Codec) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})))
	}

	pc, ok := cdc.(codec.GRPCCodecProvider)
	if !ok {
		panic("failed to cast codec to codec.GRPCCodecProvider)")
	}
	grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(grpc.ForceCodec(pc.GRPCCodec())))

	address := HTTPProtocols.ReplaceAllString(cfg.Address, "")
	return grpc.Dial(address, grpcOpts...)
}
