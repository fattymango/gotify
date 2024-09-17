package app

import (
	"fmt"
	"gotify/config"
	"gotify/internal/gotify"
	"gotify/internal/server"
	"gotify/pkg/auth"
	"gotify/pkg/bolt"
	"gotify/pkg/logger"
)

func Start() {

	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %s", err))
	}

	l, err := logger.NewLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to start logger: %s", err))
	}

	db, err := bolt.NewBoltDB(cfg)
	if err != nil {
		l.Logger.Fatalf("Failed to start db: %s", err)
	}

	// migrate.Migrate(db)

	a := auth.NewAuth(cfg)
	s := server.NewServer(cfg, l, a)

	app := gotify.NewGotify(cfg, l, db, s, a)

	err = app.Start()
	if err != nil {
		l.Logger.Fatalf("Failed to start gotify: %s", err)
	}

}
