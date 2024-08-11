package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *server) GetUser(ctx context.Context, request *user.GetUserRequest) (*user.User, error) {
	var findUser models.UserModel

	err := s.db.Where("email = ?", strings.ToLower(request.Email)).First(&findUser).Error

	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(codes.NotFound, fmt.Sprintf(
			"user with email `%s` not found",
			request.Email,
		))
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.User{
		Id:       int64(findUser.ID),
		Email:    findUser.Email,
		Password: findUser.Password,
	}, nil
}
