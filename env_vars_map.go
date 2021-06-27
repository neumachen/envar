package envar

import (
	"os"
	"strings"
)

type EnvVarsLoaderFunc func() EnvVarsMap

func loadEnvVarsToMap() EnvVarsMap {
	envStrs := os.Environ()
	if len(envStrs) < 1 {
		return nil
	}
	eMap := make(EnvVarsMap)
	for i := range envStrs {
		keyValue := strings.SplitN(envStrs[i], "=", 2)
		eMap.Set(keyValue[0], keyValue[1])
	}
	return eMap
}

type EnvVarsMap map[string]string

func (e EnvVarsMap) GetLength() int {
	return len(e)
}

func (e EnvVarsMap) Get(key string) (string, bool) {
	v, ok := e[key]
	return v, ok
}

func (e EnvVarsMap) Set(key, value string) {
	e[key] = value
}
