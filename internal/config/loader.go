package config

import (
	"os"
	"strconv"
)

type RuntimeVars struct {
	FileUrl string
	ReqFreq int // numbers of times to hit the URL in an hour
}

var runtimeConfig *RuntimeVars

func GetRuntimeConfig() *RuntimeVars {
	if runtimeConfig == nil {
		loadConfig()
	}

	return runtimeConfig
}

func loadConfig() {
	runtimeConfig = &RuntimeVars{
		FileUrl: getEnv("FILE_URL", "https://host.com/sample-file-path"),
		ReqFreq: getEnvAsInt("REQ_FREQ", 2),
	}

	// Make sure that only 60 requests are allowed in an hour, one each minute.
	if runtimeConfig.ReqFreq > 60 {
		runtimeConfig.ReqFreq = 60
	}

	// TODO: Verify URL structure using regex
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