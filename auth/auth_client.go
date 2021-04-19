package auth

import (
	"context"

	"github.com/qibobo/grpcgo/models"
	"google.golang.org/grpc"
)

type AuthClient struct {
	service  models.LoginServiceClient
	username string
	password string
}

func NewAuthClient(cc grpc.ClientConnInterface, username string, password string) *AuthClient {
	service := models.NewLoginServiceClient(cc)
	return &AuthClient{
		service:  service,
		username: username,
		password: password,
	}
}

func (ac *AuthClient) Login() (string, error) {
	resp, err := ac.service.Login(context.Background(), &models.LoginRequest{
		Username: ac.username,
		Password: ac.password,
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}
