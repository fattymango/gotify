package gotify

import (
	"context"
	"fmt"
	"gotify/model"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (g *Gotify) getAuthenticatedClient() (*spotify.Client, error) {

	authenticate := func() (*spotify.Client, error) {
		g.Logger.Logger.Info("No user found, starting authentication process")
		c, err := g.authenticate()
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate: %s", err)
		}
		g.Client = c
		return c, nil
	}
	users, err := g.GetUsers()
	if err != nil {
		g.Logger.Logger.Fatalf("Failed to get users: %s", err)
		authenticate()
	}

	if len(users) == 0 {
		return authenticate()
	}

	user := users[0]

	if user.AccessToken == "" || user.RefreshToken == "" {
		return authenticate()
	}

	g.Logger.Logger.Info("User found, refreshing token")

	token, err := g.Auth.SpotifyAuth.RefreshToken(context.Background(), &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %s", err)
	}

	c := g.newSpotifyClient(token)
	g.Client = c

	// get latest user info
	u, err := g.Client.CurrentUser(context.Background())
	if err != nil {
		xerr := g.DeleteUserBySpotifyID(user.SpotifyID)
		if xerr != nil {
			return nil, fmt.Errorf("failed to delete user: %s", err)
		}
		fmt.Println("User deleted")
		return nil, fmt.Errorf("failed to get user: %s", err)
	}

	err = g.saveToken(u, token)
	if err != nil {
		return nil, fmt.Errorf("failed to save token: %s", err)
	}

	return c, nil

}

func (g *Gotify) authenticate() (*spotify.Client, error) {

	// create a channel to receive the token
	ch := make(chan *oauth2.Token)

	// start the server
	go func() {
		err := g.Server.Start(ch)
		if err != nil {
			g.Logger.Logger.Fatalf("Failed to start server: %s", err)
		}
	}()

	// wait for the token
	token := <-ch

	// get the user
	client := g.newSpotifyClient(token)

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// save the token
	err = g.saveToken(user, token)
	if err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	g.Logger.Logger.Infof("User authenticated: %s", user.DisplayName)

	return client, nil

}

func (g *Gotify) saveToken(user *spotify.PrivateUser, token *oauth2.Token) error {
	u := model.User{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		SpotifyID:    user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
	}

	g.User = &u

	err := g.SaveUser(&u)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}
