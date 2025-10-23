package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

// WithAuth
//
// 인증 미들 웨어 입니다.
func (m *middleware) WithAuth(op huma.Operation) huma.Operation {
	op.Security = append(op.Security, map[string][]string{"bearer": {}})
	op.Middlewares = huma.Middlewares{m.authMiddleware}
	return op
}

func (m *middleware) authMiddleware(ctx huma.Context, next func(huma.Context)) {
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

	ctx = huma.WithValue(ctx, "user_id", user.ID)
	ctx = huma.WithValue(ctx, "user_role", user.Role)

	next(ctx)
}
