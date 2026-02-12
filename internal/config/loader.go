package config

import (
	"contents-api-file-monitor/internal/logger"
	"os"
	"strconv"
)

type RuntimeVars struct {
	FileUrl          string
	ReqFreq          int
	ClientTimeoutSec int
}

func LoadRuntimeConfig(log *logger.Logger) *RuntimeVars {
	cfg := &RuntimeVars{
		FileUrl:          getEnv(log, "FILE_URL", "https://host.com/sample-file-path"),
		ReqFreq:          getEnvAsInt(log, "REQ_FREQ", 2),
		ClientTimeoutSec: getEnvAsInt(log, "CLIENT_TIMEOUT_SEC", 10),
	}

	logger.Infof(log, "Initial config loaded: FileUrl: %t, ReqFreq: %t, ClientTimeoutSec: %t", cfg.FileUrl != "", cfg.ReqFreq > 0, cfg.ClientTimeoutSec > 0)

	// Make sure that only 60 requests are allowed in an hour, one each minute.
	if cfg.ReqFreq > 60 {
		cfg.ReqFreq = 60
	}

	// TODO: Verify URL structure using regex

	return cfg
}

func getEnv(log *logger.Logger, key, defVal string) string {
	logger.Infof(log, "Reading environment variable: %s", key)
	val := os.Getenv(key)
	if val == "" {
		logger.Infof(log, "Environment variable %s not set, using default: %s", key, defVal)
		return defVal
	}

	return val
}

func getEnvAsInt(log *logger.Logger, key string, defVal int) int {
	logger.Infof(log, "Reading integer environment variable: %s", key)
	val := os.Getenv(key)
	if val == "" {
		logger.Infof(log, "Environment variable %s not set, using default: %d", key, defVal)
		return defVal
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defVal
	}

	return int(i)
}