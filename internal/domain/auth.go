package domain

import (
	"context"
)

type Token struct {
	AccessToken  string `json:"access_token" doc:"엑세스 토큰입니다." example:"access_token"`
	RefreshToken string `json:"refresh_token" doc:"리프레시 토큰입니다." example:"refresh_token"`
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
