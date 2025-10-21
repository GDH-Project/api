package grpc

import (
	"context"

	"github.com/GDH-Proejct/api/cmd/config"
	"github.com/GDH-Proejct/api/internal/domain"
	"github.com/GDH-Proejct/api/internal/grpc/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userClient struct {
	log    *zap.Logger
	client userpb.UserServiceClient
}

func (c *userClient) CheckCreateUser(ctx context.Context, email string, name string) error {
	state, err := c.client.CheckCreateUser(ctx, &userpb.GetCheckCreateUserRequest{
		Email: email,
		Name:  name,
	})

	if err != nil || !state.Ok {
		err = errorFromGrpcError(err)
		c.log.Info("uc.CheckCreateUser() 실패", zap.Error(err))
		return err
	}

	return nil
}

func (c *userClient) GetUserInfoByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := c.client.GetUserInfoByEmail(ctx, &userpb.GetUserInfoByEmailRequest{
		Email: email,
	})

	if err != nil {
		err = errorFromGrpcError(err)
		c.log.Info("uc.GetUserInfoByEmail() 오류", zap.Error(err),
			zap.String("email", email),
		)
		return nil, err
	}

	u := &domain.User{
		ID:    user.GetUserId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
		Role:  parseGrpcRoleToUserRole(user.GetRole()),
	}

	return u, nil
}

func (c *userClient) GetUserInfoByUserID(ctx context.Context, id string) (*domain.User, error) {
	user, err := c.client.GetUserInfoByUserID(ctx, &userpb.GetUserInfoByUserIDRequest{
		UserId: id,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		c.log.Info("uc.GetUserInfoByUserID() 오류", zap.Error(err),
			zap.String("id", id),
		)
		return nil, err
	}

	u := &domain.User{
		ID:    user.GetUserId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
		Role:  parseGrpcRoleToUserRole(user.GetRole()),
	}

	return u, nil
}

func (c *userClient) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := c.client.CreateUser(ctx, &userpb.CreateUserRequest{
		Type:     parseUserRoleToCreateUserType(string(user.Role)),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		c.log.Info("uc.CreateUser() 오류", zap.Error(err),
			zap.String("email", user.Email),
			zap.String("name", user.Name),
			zap.String("role", string(user.Role)),
		)
		return err
	}

	return nil
}

func (c *userClient) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	userResp, err := c.client.UpdateUser(ctx, &userpb.UpdateUserRequest{
		UserId:   user.ID,
		Name:     user.Name,
		Password: user.Password,
	})
	if err != nil {
		err = errorFromGrpcError(err)
		c.log.Info("uc.UpdateUser() 오류", zap.Error(err),
			zap.String("id", user.ID),
			zap.String("email", user.Email),
			zap.String("name", user.Name),
		)
		return nil, err
	}

	return &domain.User{
		ID:    userResp.GetUserId(),
		Name:  userResp.GetName(),
		Email: userResp.GetEmail(),
		Role:  parseGrpcRoleToUserRole(userResp.GetRole()),
	}, nil
}

func (c *userClient) DeleteUser(ctx context.Context, id string, password string) error {
	_, err := c.client.DeleteUser(ctx, &userpb.DeleteUserRequest{
		UserId:   id,
		Password: password,
	})

	if err != nil {
		err = errorFromGrpcError(err)
		c.log.Warn("uc.DeleteUser() 오류", zap.Error(err),
			zap.String("id", id))
		return err
	}

	return nil
}

func NewUserClient(log *zap.Logger, cfg *config.EnvConfig) domain.UserClient {
	conn, err := grpc.NewClient(
		cfg.AuthGrpcServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("사용자 gRPC 클라이언트를 초기화 할 수 없습니다..", zap.Error(err),
			zap.String("url", cfg.AuthGrpcServer),
		)
	}

	client := userpb.NewUserServiceClient(conn)

	return &userClient{
		log:    log,
		client: client,
	}
}
