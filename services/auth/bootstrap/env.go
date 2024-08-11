package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	return &env{
		Port:               safeEnv("PORT", "8081"),
		UserServiceAddress: safeEnv("USER_SERVICE_ADDRESS", ":8080"),
		JWTSecret:          safeEnv("JWT_SECRET", ""),
	}, nil
}
