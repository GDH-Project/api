package domain

import (
	"context"
)

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleUser   UserRole = "user"
	UserRoleDevice UserRole = "device"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	Role     UserRole
}

func ParseStringRoleToUserRole(role string) UserRole {
	var r UserRole
	switch UserRole(role) {
	case UserRoleUser:
		r = UserRoleUser
	case UserRoleDevice:
		r = UserRoleDevice
	default:
		r = UserRoleUser
	}

	return r
}

type UserClient interface {
	CheckCreateUser(ctx context.Context, email string, name string) error
	GetUserInfoByEmail(ctx context.Context, email string) (*User, error)
	GetUserInfoByUserID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id string, password string) error
}

type UserService interface {
	UserClient
}

type UserUseCase interface {
	UserClient
}
