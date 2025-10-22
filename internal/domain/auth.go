package domain

import (
	"context"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthClient interface {
	Login(ctx context.Context, email string, password string) (*Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
	Logout(ctx context.Context, accessToken string) error
	Validate(ctx context.Context, accessToken string) (*User, error)
}

type AuthService interface {
	AuthClient
}

type AuthUseCase interface {
	AuthClient
}
