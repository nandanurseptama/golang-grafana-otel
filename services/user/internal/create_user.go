package internal

import (
	"context"
	"log/slog"
	"net/mail"
	"strings"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal/models"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/pkg/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *server) CreateUser(
	ctx context.Context,
	request *user.CreateUserRequest,
) (*user.User, error) {
	ctx, span := otel.Tracer.Start(ctx, "CreateUser")
	defer span.End()
	// handling request nil
	if request == nil {
		slog.ErrorContext(ctx, "cannot create user", slog.Any("reason", "request cannot be nil"))
		span.RecordError(status.Error(codes.InvalidArgument, "request cannot be nil"))
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	// handling email or password field empty
	if request.Email == "" || request.Password == "" {
		slog.ErrorContext(ctx, "email or password required")
		span.RecordError(status.Error(codes.InvalidArgument, "email and password required"))
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	// parsing email to lower case
	email := strings.ToLower(request.Email)
	// validate email address
	_, err := mail.ParseAddress(email)

	// handling when email invalid
	if err != nil {
		slog.ErrorContext(ctx, "invalid email format")
		span.RecordError(status.Error(codes.InvalidArgument, "invalid email format"))
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	var totalUser int64

	// count query with `email` from DB
	err = s.db.WithContext(ctx).Model(&models.UserModel{}).Where("email = ?", email).Count(&totalUser).Error

	// error handling when count error
	if err != nil && err != gorm.ErrRecordNotFound {
		slog.ErrorContext(ctx, "internal server error", slog.Any("reason", err.Error()))
		span.RecordError(status.Error(codes.Internal, err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check query count result
	if totalUser > 0 {
		slog.ErrorContext(ctx, "email already registered")
		span.RecordError(status.Error(codes.AlreadyExists, "emal already registered"))
		return nil, status.Error(codes.AlreadyExists, "email already registered")
	}

	// new user data
	var newUser = models.UserModel{
		Email:    request.Email,
		Password: request.Password,
	}

	// save user
	err = s.db.Save(&newUser).Error

	if err != nil {
		slog.ErrorContext(ctx, "internal server error", slog.Any("reason", err.Error()))
		span.RecordError(status.Error(codes.Internal, err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	slog.InfoContext(ctx, "success create new user")

	return &user.User{
		Id:       int64(newUser.ID),
		Email:    newUser.Email,
		Password: newUser.Password,
	}, nil
}
