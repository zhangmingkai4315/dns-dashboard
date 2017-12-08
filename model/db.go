package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	// install sqlite engine
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var globalDB *gorm.DB

func init() {
	db, err := gorm.Open("sqlite3", "sqlite.db")
	if err == nil {
		db.AutoMigrate(&DNSSerialData{})
		globalDB = db
	}
}

// GetDB return db instance
func GetDB() (*gorm.DB, error) {
	if globalDB == nil {
		return nil, errors.New("database connection fail")
	}
	return globalDB, nil
}
