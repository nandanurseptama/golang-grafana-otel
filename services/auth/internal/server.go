package internal

import (
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	user_svc "github.com/nandanurseptama/golang-grafana-otel/services/user"
)

type server struct {
	jwtSecret         string
	userServiceClient user_svc.UserServiceClient
	auth.UnimplementedAuthServiceServer
}

func NewServer(
	jwtSecret string,
	userServiceClient user_svc.UserServiceClient,
) (*server, error) {
	return &server{
		jwtSecret:         jwtSecret,
		userServiceClient: userServiceClient,
	}, nil
}
