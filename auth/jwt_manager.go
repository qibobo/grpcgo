package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qibobo/grpcgo/models"
)

type UserClaim struct {
	jwt.StandardClaims
	UserName string `json:"user_name"`
	Role     string `json:"role"`
}
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (jm *JWTManager) Generate(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserName: user.UserName,
		Role:     user.Role,
	})
	return token.SignedString([]byte(jm.secretKey))
}

func (jm *JWTManager) Verify(tokenStr string) (*UserClaim, error) {
	accessToken, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}

		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claim, ok := accessToken.Claims.(*UserClaim)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claim, nil
}
