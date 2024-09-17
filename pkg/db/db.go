package db

import (
	"gotify/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	DB     *gorm.DB
	config *config.Config
}

func NewSqlite(cfg *config.Config) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		DB: db,
	}, nil
}
