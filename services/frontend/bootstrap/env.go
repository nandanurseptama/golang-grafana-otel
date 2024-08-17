package bootstrap

import (
	"context"
	"log/slog"
	"os"
)

type env struct {
	Port               string
	AuthServiceAddress string
}

func safeEnv(key string, defaultValue string) string {
	x := os.Getenv(key)
	if x == "" {
		return defaultValue
	}
	return x
}

func GetEnv(ctx context.Context) (*env, error) {
	slog.InfoContext(ctx, "initiate env")

	return &env{
		Port:               safeEnv("PORT", "8081"),
		AuthServiceAddress: safeEnv("AUTH_SERCVICE_ADDRESS", ":8081"),
	}, nil
}
