package database

import (
	"komari/cmd/flags"
	"komari/database/models"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func InitDatabase() {
	if _, err := os.Stat(flags.DatabaseFile); os.IsNotExist(err) {
		log.Printf("Database file %q does not exist, creating...", flags.DatabaseFile)
		dbDir := filepath.Dir(flags.DatabaseFile)
		if dbDir != "" {
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				log.Fatalf("Failed to create directory %q for database file: %v", dbDir, err)
			}
		}
		file, err := os.Create(flags.DatabaseFile)
		if err != nil {
			log.Fatalf("Failed to create SQLite3 database file %q: %v", flags.DatabaseFile, err)
		}
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close database file %q: %v", flags.DatabaseFile, err)
		}
	} else if err != nil {
		log.Fatalf("Failed to check database file %q: %v", flags.DatabaseFile, err)
	}
}

func GetDBInstance() *gorm.DB {
	once.Do(func() {
		var err error
		instance, err = gorm.Open(sqlite.Open(flags.DatabaseFile), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite3 database: %v", err)
		}
		err = instance.AutoMigrate(
			&models.ClientConfig{},
			&models.User{},
			&models.Session{},
			&models.History{},
			&models.ClientInfo{},
			&models.Config{},
			&models.Custom{},
		)
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
	})
	return instance
}
