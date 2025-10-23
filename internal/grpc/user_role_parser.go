package grpc

import (
	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/grpc/userpb"
)

func parseUserRoleToCreateUserType(role string) userpb.CreateUserType {
	var r userpb.CreateUserType

	switch domain.UserRole(role) {
	case domain.UserRoleUser:
		r = userpb.CreateUserType_basic_user
	case domain.UserRoleDevice:
		r = userpb.CreateUserType_data_user
	default:
		r = userpb.CreateUserType_basic_user
	}

	return r
}

func parseGrpcRoleToUserRole(role userpb.UserRole) domain.UserRole {
	var r domain.UserRole

	switch role {
	case userpb.UserRole_BASIC_USER:
		r = domain.UserRoleUser
	case userpb.UserRole_DATA_USER:
		r = domain.UserRoleDevice
	case userpb.UserRole_ADMIN:
		r = domain.UserRoleAdmin
	default:
		r = domain.UserRoleUser
	}

	return r
}
