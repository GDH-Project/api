package config

import "go.uber.org/zap"

func InitLogger(isDebugMode bool) *zap.Logger {
	var logger *zap.Logger
	if isDebugMode {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	logger.Info("Logger Init Success")
	return logger
}
