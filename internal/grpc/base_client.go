package grpc

import (
	"github.com/GDH-Project/api/cmd/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBaseClient(log *zap.Logger, cfg *config.EnvConfig) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		cfg.AuthGrpcServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("사용자 gRPC 클라이언트를 초기화 할 수 없습니다..", zap.Error(err),
			zap.String("url", cfg.AuthGrpcServer),
		)
	}

	return conn
}
