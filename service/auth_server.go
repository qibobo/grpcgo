package service

import (
	"context"
	"log"

	"github.com/qibobo/grpcgo/auth"
	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/service/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore  store.UserStore
	jwtManager *auth.JWTManager
	models.UnimplementedLoginServiceServer
}

func NewAuthServer(userStore store.UserStore, jwtManager *auth.JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}
func (as *AuthServer) Login(ctx context.Context, loginRequest *models.LoginRequest) (*models.LoginResponse, error) {
	userName := loginRequest.GetUsername()
	password := loginRequest.GetPassword()
	log.Printf("login request username %s password %s", userName, password)
	user := as.userStore.Find(userName)
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "can not find user %s", userName)
	}
	if !user.IsCorrectPassword(password) {
		return nil, status.Error(codes.Unauthenticated, "password is not correct")
	}
	token, err := as.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}
	return &models.LoginResponse{
		AccessToken: token,
	}, nil

}
