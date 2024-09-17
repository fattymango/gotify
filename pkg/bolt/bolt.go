package bolt

import (
	"gotify/config"

	"go.etcd.io/bbolt"
)

type BoltDB struct {
	DB     *bbolt.DB
	config *config.Config
}

func NewBoltDB(cfg *config.Config) (*BoltDB, error) {
	db, err := bbolt.Open(cfg.DB.Path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltDB{
		DB: db,
	}, nil
}
