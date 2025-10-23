package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(isDebugMode bool) *zap.Logger {
	var cfg zap.Config
	if isDebugMode {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	logger, _ := cfg.Build()

	defer logger.Sync()

	logger.Info("Logger Init Success")
	return logger
}
