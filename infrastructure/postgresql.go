package infrastructure

import (
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
)

var pgClient *PGClient

type PGClient struct {
	db *gorm.DB
}

// GetDB returns the PGClient singleton.
func GetDB() *gorm.DB {
	if pgClient != nil {
		return pgClient.db
	}
	return nil
}

// InitPG initialises the postgres connection for the app.
func InitPG() error {
	dbURL := "postgres://hermes:password@localhost:5432/hermes"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Errorf("Failed to open postgres connection: %v", err)
		return err
	}
	pgClient = &PGClient{
		db: db,
	}
	return nil
}
