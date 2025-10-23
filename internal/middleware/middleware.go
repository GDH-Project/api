package middleware

import (
	"github.com/GDH-Project/api/internal/domain"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Middleware interface {
	WithAuth(op huma.Operation) huma.Operation
	WithGrpcMeta() gin.HandlerFunc
}

type middleware struct {
	api         huma.API
	log         *zap.Logger
	authUseCase domain.AuthUseCase
}

func NewMiddleware(api huma.API, log *zap.Logger, authUseCase domain.AuthUseCase) Middleware {
	return &middleware{
		api:         api,
		log:         log,
		authUseCase: authUseCase,
	}
}
