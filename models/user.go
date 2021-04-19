package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName       string
	HashedPassword string
	Role           string
}

func NewUser(userName string, password string, role string) (*User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("can not hased password %s", err)
	}
	return &User{
		UserName:       userName,
		HashedPassword: string(hashedBytes),
		Role:           role,
	}, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
