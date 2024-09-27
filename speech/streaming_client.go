package speech

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	pb "github.com/vito-ai/go-genproto/vito-openapi/stt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/vito-ai/auth"
	"github.com/vito-ai/auth/option"
)

type gRPCClient struct {
	coonPool *grpc.ClientConn
	client   pb.OnlineDecoderClient

	tp auth.TokenProvider
}

// gRPC 스트리밍을 위한 새로운 gRPC 클라이언트를 만듭니다.
func NewStreamingClient(ctx context.Context, cliopts *option.ClientOption) (*gRPCClient, error) {
	if cliopts == nil {
		cliopts = option.DefaultClientOption()
	}
	tp, err := auth.NewRTZRTokenProvider(cliopts)
	if err != nil {
		return nil, err
	}

	c := &gRPCClient{
		tp: tp,
	}

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))

	conn, err := grpc.NewClient(cliopts.GetStreamingEndpoint(), dialOpts...)
	if err != nil {
		return nil, err
	}

	if isConnReady := conn.WaitForStateChange(ctx, connectivity.Ready); !isConnReady {
		return nil, fmt.Errorf("cannot connect to rtzr stt grpc server")
	}

	client := pb.NewOnlineDecoderClient(conn)

	c.coonPool = conn
	c.client = client
	return c, nil
}

func (c *gRPCClient) StreamingRecognize(ctx context.Context) (pb.OnlineDecoder_DecodeClient, error) {
	token, err := c.tp.Token(ctx)
	if err != nil {
		return nil, err
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", "bearer", token.AccessToken))
	ctxWithAuth := metautils.NiceMD(md).ToOutgoing(ctx)

	stream, err := c.client.Decode(ctxWithAuth)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (c *gRPCClient) Close() error {
	return c.coonPool.Close()
}
