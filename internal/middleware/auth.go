package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/util"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

// WithAuth
//
// 인증 미들 웨어 입니다.
func (m *middleware) WithAuth(op huma.Operation, passRole ...domain.UserRole) huma.Operation {
	var targetRole domain.UserRole
	if len(passRole) > 0 {
		targetRole = passRole[0]
	}
	op.Security = append(op.Security, map[string][]string{"bearer": {}})
	op.Middlewares = huma.Middlewares{m.authMiddleware(targetRole)}
	return op
}

func (m *middleware) authMiddleware(passRole domain.UserRole) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			err := errors.New("인증 헤더가 없거나 유효하지 않습니다")
			m.log.Info("Authorization 헤더가 없거나 유효하지 않습니다", zap.Error(err))
			_ = huma.WriteErr(m.api, ctx, http.StatusForbidden, "Authorization 헤더가 없거나 유효하지 않습니다.", err)
			return
		}

		token := authHeader[len("Bearer "):]

		user, err := m.authUseCase.Validate(ctx.Context(), token)
		if err != nil {
			m.log.Info("accessToken이 유효하지 않습니다.", zap.Error(err),
				zap.String("token", token),
			)

			_ = huma.WriteErr(m.api, ctx, http.StatusForbidden, "accessToken이 유효하지 않습니다.", err)
			return

		}
		ctx = huma.WithValue(ctx, domain.CTX_USER_ID, user.ID)
		ctx = huma.WithValue(ctx, domain.CTX_USER_ROLE, string(user.Role))

		if passRole != "" {
			validater := util.ValidateUser{
				Ctx:        ctx.Context(),
				TargetRole: passRole,
			}

			state := validater.Exec()
			if !state {
				_ = huma.WriteErr(m.api, ctx, http.StatusForbidden, "접근 권한이 없습니다.", err)
				return
			}
		}

		next(ctx)
	}
}
