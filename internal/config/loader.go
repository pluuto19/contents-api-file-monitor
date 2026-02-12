package config

import (
	"os"
	"strconv"
)

type RuntimeVars struct {
	FileUrl          string
	ReqFreq          int
	ClientTimeoutSec int
}

func LoadRuntimeConfig() *RuntimeVars {
	cfg := &RuntimeVars{
		FileUrl:          getEnv("FILE_URL", "https://host.com/sample-file-path"),
		ReqFreq:          getEnvAsInt("REQ_FREQ", 2),
		ClientTimeoutSec: getEnvAsInt("CLIENT_TIMEOUT_SEC", 10),
	}

	// Make sure that only 60 requests are allowed in an hour, one each minute.
	if cfg.ReqFreq > 60 {
		cfg.ReqFreq = 60
	}

	// TODO: Verify URL structure using regex

	return cfg
}

func getEnv(key, defVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defVal
	}

	return val
}

func getEnvAsInt(key string, defVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defVal
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defVal
	}

	return int(i)
}