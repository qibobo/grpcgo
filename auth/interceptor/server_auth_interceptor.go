package interceptor

import (
	"context"
	"log"

	"github.com/qibobo/grpcgo/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServerInterceptor struct {
	accessibleRoles map[string][]string
	jwtManager      *auth.JWTManager
}

func NewServerInterceptor(accessibleRoles map[string][]string, jwtManager *auth.JWTManager) *ServerInterceptor {
	return &ServerInterceptor{
		accessibleRoles: accessibleRoles,
		jwtManager:      jwtManager,
	}
}
func (si *ServerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("server unary authorize interceptor %s", info.FullMethod)
		err := si.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized %s", err)
		}
		return handler(ctx, req)
	}
}
func (si *ServerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Printf("server stream authorize interceptor %s", info.FullMethod)
		err := si.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "Unauthorized %s", err)
		}
		return handler(srv, ss)
	}
}

func (si *ServerInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := si.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md[auth.AuthorizationName]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization is not provided")
	}
	log.Printf("access token is %s", values[0])
	claim, err := si.jwtManager.Verify(values[0])
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid authorization %s", err)
	}
	for _, role := range accessibleRoles {
		if role == claim.Role {
			return nil
		}
	}
	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
