package migrate

import (
	"fmt"
	"gotify/model"
	"gotify/pkg/db"
)

func Migrate(db *db.Sqlite) error {
	err := db.DB.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate user: %w", err)
	}

	return nil
}
