package database_test

import (
	"iron-stream/internal/database"
	"testing"
)

func TestGetUserByEmail(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:    "agustfricke@proton.me",
		Name:     "Agustin",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	t.Run("success", func(t *testing.T) {
		user, err := database.GetUserByEmail("agustfricke@proton.me")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if user.Email != "agustfricke@proton.me" {
			t.Errorf("expected 'agustfricke@proton.me' but got: %v", user.Email)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := database.GetUserByEmail("idontexist@email.com")
		if err == nil {
			t.Errorf("expected error but got nil")
		}
		if err.Error() != "User not found with email idontexist@email.com" {
			t.Errorf("expected error to be 'User not found with email idontexist@email.com' but got: %v", err.Error())
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users WHERE email = 'agustfricke@proton.me'`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")

	t.Run("success", func(t *testing.T) {
		input := database.User{
			Email:    "agustfricke@some.com",
			Name:     "Agustin",
			Surname:  "Fricke",
			Password: "some-password",
			Pc:       "agust@ubuntu",
		}
		err := database.CreateUser(input)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		input := database.User{
			Email:    "agustfricke@some.com",
			Name:     "Agustin",
			Surname:  "Fricke",
			Password: "some-password",
			Pc:       "agust@ubuntu",
		}
		err := database.CreateUser(input)
		if err == nil {
			t.Errorf("test failed because: %v", err)
		}
		if err.Error() != "The email: agustfricke@some.com already exists" {
			t.Errorf("expected error to be 'agustfricke@gmail.com already exists' but got: %v", err.Error())
		}
	})

	_, err := database.DB.Exec(`DELETE FROM users WHERE email = 'agustfricke@some.com'`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}
