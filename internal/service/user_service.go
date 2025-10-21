package service

import (
	"context"

	"github.com/GDH-Proejct/api/internal/domain"
	"go.uber.org/zap"
)

type userService struct {
	client domain.UserClient
	log    *zap.Logger
}

func (svc *userService) CheckCreateUser(ctx context.Context, email string, name string) error {
	return svc.client.CheckCreateUser(ctx, email, name)
}

func (svc *userService) GetUserInfoByEmail(ctx context.Context, email string) (*domain.User, error) {
	return svc.client.GetUserInfoByEmail(ctx, email)
}

func (svc *userService) GetUserInfoByUserID(ctx context.Context, id string) (*domain.User, error) {
	return svc.client.GetUserInfoByUserID(ctx, id)
}

func (svc *userService) CreateUser(ctx context.Context, user *domain.User) error {
	return svc.client.CreateUser(ctx, user)
}

func (svc *userService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return svc.client.UpdateUser(ctx, user)
}

func (svc *userService) DeleteUser(ctx context.Context, id string, password string) error {
	return svc.client.DeleteUser(ctx, id, password)
}

func NewUserService(log *zap.Logger, uc domain.UserClient) domain.UserService {
	return &userService{
		client: uc,
		log:    log,
	}
}
