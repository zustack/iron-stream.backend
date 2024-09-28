package database_test

import (
	"iron-stream/internal/database"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDeleteAccountByEmail(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:    "agustfricke@proton.me",
		Password: "some-password",
		Pc:       "agust@ubuntu",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	t.Run("success", func(t *testing.T) {
		err = database.DeleteAccountByEmail("agustfricke@proton.me")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		user, err := database.GetUserByEmail("agustfricke@proton.me")
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if user.Email != "" {
			t.Errorf("Expected email to be empty but got: %v", user.Email)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		err = database.DeleteAccountByEmail("foo@fiz.com")
		if err == nil {
			t.Errorf("Expected error got nil")
		}
		if err.Error() != "No account found with the email foo@fiz.com" {
			t.Errorf("Expected error to be 'No account found with the email foo@fiz.com' but got: %v", err.Error())
		}
	})

	t.Run("error delete user with id 1", func(t *testing.T) {
		_, err = database.DB.Exec(`DELETE FROM users;`)
		if err != nil {
			t.Fatalf("failed to teardown test database: %v", err)
		}
		database.DB.Exec(`
      DROP TABLE IF EXISTS users;
      CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
        password VARCHAR(255) NOT NULL,
        email VARCHAR(55) NOT NULL UNIQUE,
        name VARCHAR(55) NOT NULL,
        surname VARCHAR(55) NOT NULL,
        is_admin BOOL,
        special_apps BOOL DEFAULT FALSE,
        is_active BOOL DEFAULT TRUE,
        email_token INT,
        verified BOOL DEFAULT FALSE, 
        pc VARCHAR(255) DEFAULT '',  
        os VARCHAR(20) DEFAULT '',  
        created_at VARCHAR(40) NOT NULL
    );`)

		err := database.CreateUser(database.User{
			Email:    "agustfricke@proton.me",
			Password: "some-password",
			Pc:       "agust@ubuntu",
		})
		if err != nil {
			t.Errorf("test failed because of CreateUser(): %v", err)
		}

		err = database.DeleteAccountByEmail("agustfricke@proton.me")
		if err == nil {
			t.Errorf("Expected error got nil")
		}
		if err.Error() != "The account with ID 1 cannot be deleted" {
			t.Errorf("Expected error to be 'The account with ID 1 cannot be deleted' but got: %v", err.Error())
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestVerifyAccount(t *testing.T) {
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
		err = database.VerifyAccount(user.ID)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		userAfter, err := database.GetUserByEmail("agustfricke@proton.me")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if user.Verified != !userAfter.Verified {
			t.Errorf("expected %v but got: %v", !user.Verified, userAfter.Verified)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

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

	// test el hash de password
	t.Run("hash password", func(t *testing.T) {
		// get user by email
		user, err := database.GetUserByEmail("agustfricke@some.com")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("some-password"))
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
	})

	// test el hash de pc
	t.Run("hash pc", func(t *testing.T) {
		user, err := database.GetUserByEmail("agustfricke@some.com")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Pc), []byte("agust@ubuntu"))
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
	})

	_, err := database.DB.Exec(`DELETE FROM users WHERE email = 'agustfricke@some.com'`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}
