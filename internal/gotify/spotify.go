package gotify

import (
	"context"
	"fmt"
	"gotify/model"
	"io"
	"net/http"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (g *Gotify) newSpotifyClient(token *oauth2.Token) *spotify.Client {
	g.Logger.Logger.Infof("Creating new Spotify client, token: %s, expiry: %s", token.AccessToken, token.Expiry.String())
	return spotify.New(g.Server.Auth.SpotifyAuth.Client(context.Background(), token))
}

func (g *Gotify) getAllFollowers(user *model.User) error {

	// https://spclient.wg.spotify.com/user-profile-view/v3/profile/j7jtl8q682lq4gq1sqnia4xnl/followers?market=from_token
	url := fmt.Sprintf("https://spclient.wg.spotify.com/user-profile-view/v3/profile/%s/followers?market=from_token", user.SpotifyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get followers: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	fmt.Println(string(body))
	return nil

}

func (g *Gotify) getLikedSongs() ([]spotify.SavedTrack, error) {

	likedSongsFetchWorker := func() ([]spotify.SavedTrack, error) {
		g.Logger.Logger.Info("Fetching liked songs worker started")
		var tracks []spotify.SavedTrack

		// Get the user's saved tracks
		ctx := context.Background()

		// Fetch the first 50 saved tracks (you can paginate through the rest)
		limit := 50
		offset := 0
		for {

			// Use the client to get user's saved tracks
			t, err := g.Client.CurrentUsersTracks(ctx, spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return nil, fmt.Errorf("error while collecting saved tracks, err = %s", err.Error())
			}

			tracks = append(tracks, t.Tracks...)

			// If no more tracks, break out of the loop
			if len(t.Tracks) < limit {
				break
			}

			// Increment offset for next batch
			offset += limit
		}

		// Save the tracks to the database
		err := g.SaveTracksForUser(g.User, tracks)
		if err != nil {
			return nil, fmt.Errorf("error while saving tracks, err = %s", err.Error())
		}

		g.Logger.Logger.Infof("Worker fetched %d liked songs", len(tracks))

		return tracks, nil

	}

	// try to get tracks from db
	tracks, err := g.GetTracksForUser(g.User)
	if err == nil { // no error, return the cached data

		// return the cached data and start a goroutine to fetch new data
		go func() {
			_, err := likedSongsFetchWorker()
			if err != nil {
				g.Logger.Logger.Errorf("error while fetching liked songs: %s", err)
			}

		}()
		return tracks, nil
	}

	// error while fetching from db, fetch from spotify
	tracks, err = likedSongsFetchWorker()
	if err != nil {
		return nil, fmt.Errorf("error while fetching liked songs: %s", err)
	}

	return tracks, nil
}
