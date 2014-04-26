package slimgfast

import (
	"os"
	"strconv"
	"strings"
)

func environ() map[string]string {
	_env := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		_env[splits[0]] = splits[1]
	}
	return _env
}

func GetEnvString(key, def string) string {
	resp, ok := environ()[key]
	if !ok {
		return def
	}
	return resp
}

func GetEnvInt(key string, def int) int {
	rawVal, ok := environ()[key]
	if !ok {
		return def
	}
	resp, err := strconv.Atoi(rawVal)
	if err != nil {
		return def
	}
	return resp
}
