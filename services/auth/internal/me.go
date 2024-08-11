package internal

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/internal/models"
	user_svc "github.com/nandanurseptama/golang-grafana-otel/services/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Me(ctx context.Context, request *auth.MeRequest) (*user_svc.User, error) {
	token, err := jwt.ParseWithClaims(
		request.Token,
		&models.JwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse auth token")
	}
	claims, ok := token.Claims.(*models.JwtClaim)

	if !ok {
		return nil, status.Error(codes.Internal, "unknwon jwt claim")
	}
	data := claims.Data
	if data == nil {
		return nil, status.Error(codes.Internal, "claim data nil")
	}

	email, ok := data["email"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "claim data nil")
	}

	res, err := s.userServiceClient.GetUser(ctx, &user_svc.GetUserRequest{
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	return res, nil

}
