package auth

import (
	"gotify/config"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Auth struct {
	config      *config.Config
	SpotifyAuth *spotifyauth.Authenticator
}

var Scopes = []string{
	spotifyauth.ScopePlaylistModifyPrivate,
	spotifyauth.ScopePlaylistModifyPublic,
	spotifyauth.ScopePlaylistReadPrivate,
	spotifyauth.ScopePlaylistReadCollaborative,
	spotifyauth.ScopeUserFollowModify,
	spotifyauth.ScopeUserFollowRead,
	spotifyauth.ScopeUserLibraryModify,
	spotifyauth.ScopeUserLibraryRead,
	spotifyauth.ScopeUserModifyPlaybackState,
	spotifyauth.ScopeUserReadCurrentlyPlaying,
	spotifyauth.ScopeUserReadPlaybackState,
	spotifyauth.ScopeUserReadPrivate,
	spotifyauth.ScopeUserReadRecentlyPlayed,
	spotifyauth.ScopeUserTopRead,
	spotifyauth.ScopeUserReadEmail,
	spotifyauth.ScopeStreaming,
}

func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		config: cfg,
		SpotifyAuth: spotifyauth.New(
			spotifyauth.WithRedirectURL(cfg.Spotify.RedirectURL),
			spotifyauth.WithClientID(cfg.Spotify.ClientID),
			spotifyauth.WithClientSecret(cfg.Spotify.ClientSecret),
			spotifyauth.WithScopes(Scopes...),
		),
	}
}
