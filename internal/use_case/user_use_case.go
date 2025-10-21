package usecase

import (
	"context"

	"github.com/GDH-Proejct/api/internal/domain"
	"go.uber.org/zap"
)

type userUseCase struct {
	svc domain.UserService
	log *zap.Logger
}

func (uc *userUseCase) CheckCreateUser(ctx context.Context, email string, name string) error {
	return uc.svc.CheckCreateUser(ctx, email, name)
}

func (uc *userUseCase) GetUserInfoByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.svc.GetUserInfoByEmail(ctx, email)
}

func (uc *userUseCase) GetUserInfoByUserID(ctx context.Context, id string) (*domain.User, error) {
	return uc.svc.GetUserInfoByUserID(ctx, id)
}

func (uc *userUseCase) CreateUser(ctx context.Context, user *domain.User) error {
	return uc.svc.CreateUser(ctx, user)
}

func (uc *userUseCase) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return uc.svc.UpdateUser(ctx, user)
}

func (uc *userUseCase) DeleteUser(ctx context.Context, id string, password string) error {
	return uc.svc.DeleteUser(ctx, id, password)
}

func NewUserUseCase(logger *zap.Logger, svc domain.UserService) domain.UserUseCase {
	return &userUseCase{
		svc: svc,
		log: logger,
	}
}
