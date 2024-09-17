package model

type User struct {
	SpotifyID    string `gorm:"primaryKey;autoIncrement:true;column:spotify_id" json:"spotify_id"`
	Email        string `gorm:"column:email" json:"email"`
	DisplayName  string `gorm:"column:display_name" json:"display_name"`
	AccessToken  string `gorm:"column:access_token" json:"access_token"`
	RefreshToken string `gorm:"column:refresh_token" json:"refresh_token"`
	RetrievedAt  int64  `gorm:"column:retrieved_at" json:"retrieved_at"`
}
