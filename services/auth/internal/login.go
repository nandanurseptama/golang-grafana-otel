package internal

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth"
	"github.com/nandanurseptama/golang-grafana-otel/services/auth/internal/models"
	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {

	// handling request nil
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	// handling email or password field empty
	if request.Email == "" || request.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	// parsing email to lower case
	email := strings.ToLower(request.Email)
	// validate email address
	_, err := mail.ParseAddress(email)

	// handling when email invalid
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	r, err := s.userServiceClient.GetUser(ctx, &user.GetUserRequest{
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	if r.Password != request.Password {
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
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to signed jwt. err : %s", err.Error()))
	}
	return &auth.LoginResponse{
		Token: jwtToken,
	}, nil
}
