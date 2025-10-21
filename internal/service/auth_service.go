package service

import (
	"context"

	"github.com/GDH-Proejct/api/internal/domain"
	"go.uber.org/zap"
)

type authService struct {
	authClient domain.AuthClient
	log        *zap.Logger
}

func (svc *authService) Login(ctx context.Context, email string, password string) (*domain.Token, error) {
	return svc.authClient.Login(ctx, email, password)
}

func (svc *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	return svc.authClient.RefreshToken(ctx, refreshToken)
}

func (svc *authService) Logout(ctx context.Context, accessToken string) error {
	return svc.authClient.Logout(ctx, accessToken)
}

func (svc *authService) Validate(ctx context.Context, accessToken string) (*domain.User, error) {
	return svc.authClient.Validate(ctx, accessToken)
}

func NewAuthService(logger *zap.Logger, client domain.AuthClient) domain.AuthService {
	return &authService{
		authClient: client,
		log:        logger,
	}
}
