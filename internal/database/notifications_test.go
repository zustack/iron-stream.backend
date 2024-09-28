package database_test

import (
	"iron-stream/internal/database"
	"testing"
)

func TestCreateNotification(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")

	t.Run("success", func(t *testing.T) {
		err := database.CreateNotification("user", "test@email.com")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
	})

	_, err := database.DB.Exec(`DELETE FROM notifications`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}
