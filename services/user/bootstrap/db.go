package bootstrap

import (
	"context"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open sqlite database from `dbPath`
func OpenDB(ctx context.Context, dbPath string) (*gorm.DB, error) {
	slog.Info("opening database", slog.Any("path", dbPath))
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
