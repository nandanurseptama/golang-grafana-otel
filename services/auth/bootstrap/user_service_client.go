package bootstrap

import (
	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserServiceClient(address string) (user.UserServiceClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithCredentialsBundle(
			insecure.NewBundle(),
		),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		return nil, err
	}

	return user.NewUserServiceClient(conn), nil
}
