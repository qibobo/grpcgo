package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/qibobo/grpcgo/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor struct {
	authClient      *auth.AuthClient
	authMethods     map[string]bool
	accessToken     string
	refreshDuration time.Duration
}

func NewClientInterceptor(authClient *auth.AuthClient, authMethods map[string]bool, refreshDuration time.Duration) (*ClientInterceptor, error) {
	c := ClientInterceptor{
		authClient:      authClient,
		authMethods:     authMethods,
		refreshDuration: refreshDuration,
	}
	err := c.RefreshTokenByInterval()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (ci *ClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("client unary authorize interceptor %s", method)
		if ci.authMethods[method] {
			ctx = ci.AttachToken(ctx)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
func (ci *ClientInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Printf("client stream authorize interceptor %s", method)
		if ci.authMethods[method] {
			log.Printf("client stream authorize interceptor attach token")
			ctx = ci.AttachToken(ctx)
		}
		return streamer(ctx, desc, cc, method, opts...)

	}
}

func (ci *ClientInterceptor) AttachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, auth.AuthorizationName, ci.accessToken)
}

func (ci *ClientInterceptor) RefreshTokenByInterval() error {
	err := ci.refreshToken()
	if err != nil {
		return err
	}
	go func() {
		wait := ci.refreshDuration
		for {
			time.Sleep(wait)
			err := ci.refreshToken()
			if err != nil {
				wait = time.Second
			}
			wait = ci.refreshDuration
		}
	}()
	return nil
}
func (ci *ClientInterceptor) refreshToken() error {
	token, err := ci.authClient.Login()
	if err != nil {
		return err
	}
	ci.accessToken = token
	return nil
}
