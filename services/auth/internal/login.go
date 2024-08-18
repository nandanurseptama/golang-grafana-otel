package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/internal/models"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/pkg/otel"
	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "Login")
	defer span.End()
	// handling request nil
	if request == nil {
		span.RecordError(status.Error(codes.InvalidArgument, "request cannot be nil"))
		slog.ErrorContext(ctx, "failed login user", slog.Any("reason", "request nil"))
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	// handling email or password field empty
	if request.Email == "" || request.Password == "" {
		span.RecordError(status.Error(codes.InvalidArgument, "email and password required"))
		slog.ErrorContext(ctx, "failed login user", slog.Any("reason", "email and password required"))
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	// parsing email to lower case
	email := strings.ToLower(request.Email)
	// validate email address
	_, err := mail.ParseAddress(email)

	// handling when email invalid
	if err != nil {
		span.RecordError(status.Error(codes.InvalidArgument, "invalid email format"))
		slog.ErrorContext(ctx, "failed login user", slog.Any("reason", "invalid email format"))
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	r, err := s.userServiceClient.GetUser(ctx, &user.GetUserRequest{
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	if r.Password != request.Password {
		span.RecordError(status.Error(codes.InvalidArgument, "email or password not same"))
		slog.ErrorContext(ctx, "failed login user", slog.Any("reason", "email or password not same"))
		return nil, status.Error(codes.InvalidArgument, "email or password not same")
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth",
			Subject:   "auth_token",
			Audience:  []string{"all"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
		Data: map[string]any{
			"id":    r.Id,
			"email": r.Email,
		},
	}).SignedString([]byte(s.jwtSecret))

	if err != nil {
		span.RecordError(status.Error(codes.InvalidArgument, err.Error()))
		slog.ErrorContext(ctx, "jwt parsing error", slog.Any("reason", err.Error()))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to signed jwt. err : %s", err.Error()))
	}
	slog.InfoContext(ctx, "success signing user")
	return &auth.LoginResponse{
		Token: jwtToken,
	}, nil
}
