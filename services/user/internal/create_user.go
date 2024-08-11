package internal

import (
	"context"
	"net/mail"
	"strings"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *server) CreateUser(
	ctx context.Context,
	request *user.CreateUserRequest,
) (*user.User, error) {
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

	var totalUser int64

	// count query with `email` from DB
	err = s.db.Model(&models.UserModel{}).Where("email = ?", email).Count(&totalUser).Error

	// error handling when count error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check query count result
	if totalUser > 0 {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.User{
		Id:       int64(newUser.ID),
		Email:    newUser.Email,
		Password: newUser.Password,
	}, nil
}
