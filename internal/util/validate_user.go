package util

import (
	"context"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/middleware"
)

// ValidateUser 유저 권한 검사
//
// auth middleware를 통과한 ctx(Ctx)를 넘기고 통과시킬 권한명(TargetRole)과 차단 권한명(ExceptionRole) 을 필요에 따라 설정하면 된다.
// 권한 평가의 우선 순위는 현재 권한을 불라 올수 없는 경우 > 차단 권한 > 통과 권한 의 우선순위를 가지고 있다.
type ValidateUser struct {
	Ctx           context.Context
	TargetRole    domain.UserRole
	ExceptionRole domain.UserRole
}

func (c *ValidateUser) getValue(key string) string {
	v := c.Ctx.Value(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

// UserID ctx에서 추출한 UserID 이다. 만약 ctx에 설전정 UserID가 없을 경우 공백 문자를 반환한다.
func (c *ValidateUser) UserID() string {
	return c.getValue(middleware.CTX_USER_ID)
}

// UserRole ctx에서 추출한 UserRole 이다. 만약 ctx에 설정된 UserRole이 없을 경우 공백 문자를 반환한다.
func (c *ValidateUser) UserRole() string {
	return c.getValue(middleware.CTX_USER_ROLE)
}

// Exec 권한 평가 실행한 후 불리언을 반환한다.
//
// 권한 평가의 우선 순위는 현재 권한을 불라 올수 없는 경우 > 차단 권한 > 통과 권한 의 우선순위를 가지고 있다.
func (c *ValidateUser) Exec() bool {
	currentRole := c.UserRole()

	if currentRole == "" || currentRole == string(c.ExceptionRole) {
		return false
	}

	if currentRole != string(c.TargetRole) {
		return false
	}

	return true

}
