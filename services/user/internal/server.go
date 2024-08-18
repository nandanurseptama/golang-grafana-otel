package internal

import (
	"fmt"

	"github.com/nandanurseptama/golang-grafana-otel/services/user"
	"github.com/nandanurseptama/golang-grafana-otel/services/user/internal/models"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
	user.UnimplementedUserServiceServer
}

func NewServer(
	db *gorm.DB,
) (*server, error) {

	err := db.AutoMigrate(&models.UserModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate `users` table : %s", err.Error())
	}

	// seeding user
	db.Save(&models.UserModel{
		Email:    "doe@gmail.com",
		Password: "12345",
	})
	// seeding user
	db.Save(&models.UserModel{
		Email:    "john@gmail.com",
		Password: "12345",
	})

	return &server{
		db: db,
	}, nil
}
