package version

import (
	_ "embed"
	"strings"

	"go.uber.org/zap"
)

//go:embed version.txt
var versionContent string

func GetVersion(log *zap.Logger) string {
	version := strings.TrimSpace(versionContent)

	log.Info("Version", zap.String("version", version))
	return version
}
