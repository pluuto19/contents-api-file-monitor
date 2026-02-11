package config

import (
	"os"
	"strconv"
)

type RuntimeVars struct {
	fileUrl string
	reqFreq int // numbers of times to hit the URL in an hour
}

var runtimeConfig *RuntimeVars

func LoadConfig() {
	if runtimeConfig == nil {
		runtimeConfig = &RuntimeVars{
			fileUrl: getEnv("FILE_URL", "https://host.com/sample-file-path"),
			reqFreq: getEnvAsInt("REQ_FREQ", 2),
		}
	}
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