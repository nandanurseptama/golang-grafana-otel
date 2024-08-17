package bootstrap

import (
	"context"
	"log/slog"

	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AuthServiceClient(ctx context.Context, address string) (auth.AuthServiceClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithCredentialsBundle(
			insecure.NewBundle(),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
	)

	if err != nil {
		slog.ErrorContext(ctx, "failed to initiate auth service client", slog.Any("reason", err.Error()))
		return nil, err
	}

	return auth.NewAuthServiceClient(conn), nil
}
