package bootstrap

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open sqlite database from `dbPath`
func OpenDB(dbPath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
