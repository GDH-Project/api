package middleware

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// WithGrpcMeta grpc 메타 데이터 설정
//
// 해당 미들웨어는 gin.Use()로 사용을 해야 한다.
func (m *middleware) WithGrpcMeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		ctx := c.Request.Context()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-user-agent", userAgent)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-ip", clientIP)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
