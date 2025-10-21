package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type EnvConfig struct {
	DbUrl          string `env:"DB_URL,required"`
	AuthGrpcServer string `env:"AUTH_GRPC_SERVER,required"`
	HostUrl        string `env:"HOST_URL"`
}

func GetConfig(log *zap.Logger) *EnvConfig {
	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatal("env 로드에 실패했습니다.", zap.Error(err))
	}

	log.Info("환경변수 로드 성공")

	return config
}
