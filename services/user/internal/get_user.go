package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal/models"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/pkg/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *server) GetUser(ctx context.Context, request *user.GetUserRequest) (*user.User, error) {
	ctx, span := otel.Tracer.Start(ctx, "GetUser")
	defer span.End()
	var findUser models.UserModel

	err := s.db.WithContext(ctx).Where("email = ?", strings.ToLower(request.Email)).First(&findUser).Error

	if err == gorm.ErrRecordNotFound {
		slog.ErrorContext(ctx, "email not found")
		span.RecordError(
			status.Error(codes.NotFound, "email not found"),
		)
		return nil, status.Error(codes.NotFound, fmt.Sprintf(
			"user with email `%s` not found",
			request.Email,
		))
	}

	if err != nil {
		slog.ErrorContext(ctx, "email not found")
		span.RecordError(
			status.Error(codes.Internal, err.Error()),
		)
		return nil, status.Error(codes.Internal, err.Error())
	}
	slog.InfoContext(ctx, "success get user")

	return &user.User{
		Id:       int64(findUser.ID),
		Email:    findUser.Email,
		Password: findUser.Password,
	}, nil
}
