package gotify

import (
	"encoding/json"
	"fmt"
	"gotify/model"
	"strings"

	"github.com/zmb3/spotify/v2"
	bolt "go.etcd.io/bbolt"
)

var (
	SpotifyBucket = "spotify_bucket"

	formatUserKey = func(id string) string {
		return fmt.Sprintf("user_%s", id)
	}
	formatUserTracksKey = func(id string) string {
		return fmt.Sprintf("tracks_%s", id)
	}
)

func (g *Gotify) GetUsers() ([]model.User, error) {
	var users []model.User

	err := g.DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, no users")
		}

		return b.ForEach(func(k, v []byte) error {
			key := string(k)
			// Check if the key matches the user_* pattern
			if strings.HasPrefix(key, "user_") {
				var user model.User
				err := json.Unmarshal(v, &user)
				if err != nil {
					return fmt.Errorf("failed to unmarshal user: %w", err)
				}
				users = append(users, user)
			}
			return nil
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil

}

func (g *Gotify) SaveUser(user *model.User) error {

	g.Logger.Logger.Infof("Saving user: %v", user)
	// check if user exists in db and update, else create
	err := g.DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			// create bucket
			b, _ = tx.CreateBucket([]byte(SpotifyBucket))
		}

		// marshal user
		userBytes, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("failed to marshal user: %w", err)
		}
		// update user
		return b.Put([]byte(formatUserKey(user.SpotifyID)), userBytes)

	})

	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}

	return nil

}

func (g *Gotify) GetUserBySpotifyID(spotify_id string) (*model.User, error) {
	user := &model.User{}
	err := g.DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, no users")
		}

		userBytes := b.Get([]byte(formatUserKey(spotify_id)))
		if userBytes == nil {
			return fmt.Errorf("user not found")
		}

		return json.Unmarshal(userBytes, user)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (g *Gotify) DeleteUserBySpotifyID(spotify_id string) error {

	err := g.DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, no users")
		}

		return b.Delete([]byte(formatUserKey(spotify_id)))
	})

	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	return nil
}

func (g *Gotify) SaveTracksForUser(user *model.User, tracks []spotify.SavedTrack) error {

	err := g.DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			// create bucket
			b, _ = tx.CreateBucket([]byte(SpotifyBucket))

		}

		// marshal tracks
		tracksBytes, err := json.Marshal(tracks)
		if err != nil {
			return fmt.Errorf("failed to marshal tracks: %w", err)
		}
		// update user
		return b.Put([]byte(formatUserTracksKey(user.SpotifyID)), tracksBytes)

	})

	if err != nil {
		return fmt.Errorf("failed to save tracks for user: %w", err)
	}

	return nil

}

func (g *Gotify) GetTracksForUser(user *model.User) ([]spotify.SavedTrack, error) {
	var tracks []spotify.SavedTrack

	err := g.DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SpotifyBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, no users")
		}

		tracksBytes := b.Get([]byte(formatUserTracksKey(user.SpotifyID)))
		if tracksBytes == nil {
			return fmt.Errorf("tracks not found")
		}

		return json.Unmarshal(tracksBytes, &tracks)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get tracks for user: %w", err)
	}

	return tracks, nil
}
