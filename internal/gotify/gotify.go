package gotify

import (
	"context"
	"fmt"
	"gotify/config"
	"gotify/internal/server"
	"gotify/model"
	"gotify/pkg/auth"
	"gotify/pkg/bolt"
	"gotify/pkg/logger"

	"github.com/zmb3/spotify/v2"
)

type Gotify struct {
	DB     *bolt.BoltDB
	Config *config.Config
	Logger *logger.Logger
	Server *server.Server
	Client *spotify.Client
	Auth   *auth.Auth
	User   *model.User
}

func NewGotify(cfg *config.Config, logger *logger.Logger, db *bolt.BoltDB, s *server.Server, a *auth.Auth) *Gotify {
	return &Gotify{
		DB:     db,
		Config: cfg,
		Logger: logger,
		Server: s,
		Auth:   a,
	}
}

func (g *Gotify) Start() error {

	// get the authenticated client
	_, err := g.getAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("failed to get authenticated client: %s", err)
	}

	res, err := g.Client.CurrentUser(context.Background())
	if err != nil {
		g.Logger.Logger.Fatalf("Failed to get user: %s", err)
	}

	fmt.Println("User followrs: ", res.Followers.Count)

	tracks, err := g.getLikedSongs()
	if err != nil {
		return err
	}

	g.Logger.Logger.Infof("tracks length %d", len(tracks))

	user, err := g.GetUserBySpotifyID(res.ID)
	if err != nil {
		return fmt.Errorf("failed to get user by spotify id: %s", err)
	}

	g.getAllFollowers(user)
	select {}

	return nil
}
