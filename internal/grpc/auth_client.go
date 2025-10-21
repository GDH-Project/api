package grpc

import (
	"context"

	"github.com/GDH-Proejct/api/internal/domain"
	"github.com/GDH-Proejct/api/internal/grpc/authpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type authClient struct {
	log    *zap.Logger
	client authpb.AuthServiceClient
}

func (ac *authClient) Login(ctx context.Context, email string, password string) (*domain.Token, error) {
	token, err := ac.client.Login(ctx, &authpb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		ac.log.Debug("ac.Login() 실패", zap.Error(err),
			zap.String("email", email),
		)
		return nil, err
	}

	return &domain.Token{
		AccessToken:  token.GetAccessToken(),
		RefreshToken: token.GetRefreshToken(),
	}, nil

}

func (ac *authClient) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	token, err := ac.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		ac.log.Debug("ac.RefreshToken() 실패", zap.Error(err),
			zap.String("refreshToken", refreshToken),
		)
		return nil, err
	}

	return &domain.Token{
		AccessToken:  token.GetAccessToken(),
		RefreshToken: token.GetRefreshToken(),
	}, nil
}

func (ac *authClient) Logout(ctx context.Context, accessToken string) error {
	_, err := ac.client.Logout(ctx, &authpb.LogoutRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		ac.log.Debug("ac.Logout() 실패", zap.Error(err),
			zap.String("accessToken", accessToken),
		)

		return err
	}

	return nil
}

func (ac *authClient) Validate(ctx context.Context, accessToken string) (*domain.User, error) {
	user, err := ac.client.Validate(ctx, &authpb.ValidateRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		ac.log.Debug("ac.Validate() 실패", zap.Error(err),
			zap.String("accessToken", accessToken),
		)
		return nil, err
	}

	var userRole domain.UserRole

	switch user.GetUserRole() {
	case authpb.UserRole_BASIC_USER:
		userRole = domain.UserRoleUser
	case authpb.UserRole_DATA_USER:
		userRole = domain.UserRoleDevice
	case authpb.UserRole_ADMIN:
		userRole = domain.UserRoleAdmin
	}

	return &domain.User{
		ID:   user.GetUserId(),
		Role: userRole,
	}, nil

}

func NewAuthClient(logger *zap.Logger, conn *grpc.ClientConn) domain.AuthClient {
	client := authpb.NewAuthServiceClient(conn)

	return &authClient{
		log:    logger,
		client: client,
	}
}
