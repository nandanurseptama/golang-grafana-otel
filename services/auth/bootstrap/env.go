package bootstrap

import (
	"os"
)

type env struct {
	Port               string
	UserServiceAddress string
	JWTSecret          string
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
		Port:               safeEnv("PORT", "8081"),
		UserServiceAddress: safeEnv("USER_SERVICE_ADDRESS", ":8080"),
		JWTSecret:          safeEnv("JWT_SECRET", "B613679A0814D9EC772F95D778C35FC5FF1697C493715653C6C712144292C5AD"),
	}, nil
}
