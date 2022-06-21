package infrastructure

import (
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

var pgClient *PGClient

type PGClient struct {
	db *gorm.DB
}

// Returns the PGClient singleton.
func GetDB() *gorm.DB {
	if pgClient != nil {
		return pgClient.db
	}
	return nil
}

// InitPG initializes the postgres connection for the app.
func InitPG() error {
	dbURL := "postgres://hermes:password@localhost:5432/hermes"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return err
	}
	pgClient = &PGClient{
		db: db,
	}
	return nil
}
