package config

import (
	"bufio"
	"os"

	"go.uber.org/zap"
)

func GetVersion(log *zap.Logger) string {
	var version string
	f, err := os.Open("version.txt")
	if err != nil {
		log.Fatal("Failed to open version.txt file", zap.Error(err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		version = scanner.Text()
	} else {
		log.Fatal("Failed to read version.txt file", zap.Error(scanner.Err()))
	}

	log.Info("Version", zap.String("version", version))

	return version
}
