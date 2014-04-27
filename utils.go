package slimgfast

import (
	"os"
	"strconv"
	"strings"
)

// environ builds a full mapping of environment variables
func environ() map[string]string {
	_env := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		_env[splits[0]] = splits[1]
	}
	return _env
}

// GetEnvString tries first to get a string from the environment, but falls
// back on a default provided value.
func GetEnvString(key, def string) string {
	resp, ok := environ()[key]
	if !ok {
		return def
	}
	return resp
}

// GetEnvInt tries first to get and parse an int from the environment, but
// falls back on a default provided value.
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
