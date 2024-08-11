package bootstrap

import (
	"os"
)

type env struct {
	Port   string
	DBPath string
}

func safeEnv(key string, defaultValue string) string {
	x := os.Getenv(key)
	if x == "" {
		return defaultValue
	}
	return x
}

func GetEnv() (*env, error) {

	return &env{
		Port:   safeEnv("PORT", "8080"),
		DBPath: safeEnv("DB_PATH", "../../volumes/services/user/sqlite.db"),
	}, nil
}
