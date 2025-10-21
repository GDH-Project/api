package usecase

import (
	"context"

	"github.com/GDH-Proejct/api/internal/domain"
	"go.uber.org/zap"
)

type authUseCase struct {
	authService domain.AuthService
	log         *zap.Logger
}

func (uc *authUseCase) Login(ctx context.Context, email string, password string) (*domain.Token, error) {
	token, err := uc.authService.Login(ctx, email, password)
	if err != nil {
		uc.log.Debug("auc.Login() 실패", zap.Error(err),
			zap.String("email", email),
		)
		return nil, err
	}
	return token, nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	token, err := uc.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		uc.log.Debug("auc.RefreshToken() 실패", zap.Error(err))
		return nil, err
	}
	return token, nil
}

func (uc *authUseCase) Logout(ctx context.Context, accessToken string) error {
	err := uc.authService.Logout(ctx, accessToken)
	if err != nil {
		uc.log.Debug("auc.Logout() 실패", zap.Error(err))
		return err
	}

	return nil
}

// Validate
// *domain.User 는 ID 와 Role 만 리턴함
// 사용시 유의할것
func (uc *authUseCase) Validate(ctx context.Context, accessToken string) (*domain.User, error) {
	// Validate로는 유저의 ID만 존재한다.
	user, err := uc.authService.Validate(ctx, accessToken)
	if err != nil {
		uc.log.Debug("auc.Validate() 실패", zap.Error(err),
			zap.String("msg", "accessToken이 유효하지 않습니다."),
		)
		return nil, err
	}

	// user 정보를 채워야 하는지는 확인 필요

	return user, nil
}

func NewAuthService(logger *zap.Logger, authService domain.AuthService) domain.AuthUseCase {
	return &authUseCase{
		authService: authService,
		log:         logger,
	}
}
