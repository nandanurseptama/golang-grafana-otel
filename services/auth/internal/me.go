package internal

import (
	"context"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/internal/models"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/otel"
	user_svc "github.com/nandanurseptama/golang-grafana-otel/services/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Me(ctx context.Context, request *auth.MeRequest) (*user_svc.User, error) {
	ctx, span := otel.Tracer.Start(ctx, "Me")
	defer span.End()
	token, err := jwt.ParseWithClaims(
		request.Token,
		&models.JwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		})

	if err != nil {
		span.RecordError(status.Error(codes.Internal, err.Error()))
		slog.ErrorContext(ctx, "failed parse jwt token", slog.Any("reason", err.Error()))
		return nil, status.Error(codes.Internal, "failed to parse auth token")
	}
	claims, ok := token.Claims.(*models.JwtClaim)

	if !ok {
		span.RecordError(status.Error(codes.Internal, "unknown jwt claim"))
		slog.ErrorContext(ctx, "failed parse jwt token", slog.Any("reason", "unknown jwt claim"))
		return nil, status.Error(codes.Internal, "unknwon jwt claim")
	}
	data := claims.Data
	if data == nil {
		span.RecordError(status.Error(codes.Internal, "claim data nil"))
		slog.ErrorContext(ctx, "failed parse jwt token", slog.Any("reason", "claim data nil"))
		return nil, status.Error(codes.Internal, "claim data nil")
	}

	email, ok := data["email"].(string)
	if !ok {
		span.RecordError(status.Error(codes.Internal, "claim email nil"))
		slog.ErrorContext(ctx, "failed parse jwt token", slog.Any("reason", "claim email nil"))
		return nil, status.Error(codes.Internal, "claim data nil")
	}

	res, err := s.userServiceClient.GetUser(ctx, &user_svc.GetUserRequest{
		Email: email,
	})

	if err != nil {
		span.RecordError(status.Error(codes.Internal, err.Error()))
		slog.ErrorContext(ctx, "failed to get user", slog.Any("reason", err.Error()))
		return nil, err
	}

	slog.InfoContext(ctx, "success validate user")
	return res, nil

}
