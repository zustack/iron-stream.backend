package database_test

import (
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGetUserCount(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "agustfricke@gmail.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "agustfricke@some.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Mac",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	now := utils.FormattedDate()
	count, err := database.GetUserCount(now, "Linux")
	if err != nil {
		t.Errorf("test failed because of GetUserCount(): %v", err)
		return
	}
	if count != 2 {
		t.Errorf("expected count to be 2 but got %d", count)
	}
	count, err = database.GetUserCount(now, "Mac")
	if err != nil {
		t.Errorf("test failed because of GetUserCount(): %v", err)
		return
	}
	if count != 1 {
		t.Errorf("expected count to be 1 but got %d", count)
	}

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestUpdateEmailToken(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("user not found", func(t *testing.T) {
		err = database.UpdateEmailToken("fo@fo.com", 794209)
		if err.Error() != "No account found with the email fo@fo.com" {
			t.Errorf("Expected error to be 'No account found with the email fo@fo.com' but got %v", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		err = database.UpdateEmailToken("agustfricke@proton.me", 794209)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err := database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.EmailToken != 794209 {
			t.Errorf("expected email token to be 794209 but got %v", user.EmailToken)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestGetUserByID(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("user not found", func(t *testing.T) {
		_, err := database.GetUserByID("99999")
		if err.Error() != "No account found with the id 99999" {
			t.Errorf("Expected error to be 'No account found with the id 99999' but got %v", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		user, err := database.GetUserByID("1")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		if user.ID != 1 {
			t.Errorf("expected id to be 1 but got %v", user.ID)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}

func TestUpdateAdminStatus(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("user not found", func(t *testing.T) {
		err = database.UpdateAdminStatus("99999", "true")
		if err.Error() != "No account found with the id 99999" {
			t.Errorf("Expected error to be 'No account found with the id 99999' but got %v", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		err = database.UpdateAdminStatus("1", "true")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err := database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.IsAdmin != true {
			t.Errorf("expected IsAdmin to be true but got %v", user.IsAdmin)
		}

		err = database.UpdateAdminStatus("1", "false")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err = database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.IsAdmin != false {
			t.Errorf("expected IsAdmin to be false but got %v", user.IsAdmin)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestUpdateUserSpecialApps(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("user not found", func(t *testing.T) {
		err = database.UpdateUserSpecialApps("99999", "true")
		if err.Error() != "No account found with the id 99999" {
			t.Errorf("Expected error to be 'No account found with the id 99999' but got %v", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		err = database.UpdateUserSpecialApps("1", "true")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err := database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.SpecialApps != true {
			t.Errorf("expected SpecialApps to be true but got %v", user.SpecialApps)
		}

		err = database.UpdateUserSpecialApps("1", "false")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err = database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.SpecialApps != false {
			t.Errorf("expected SpecialApps to be false but got %v", user.SpecialApps)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestUpdateActiveStatusAllUsers(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "agustfricke@gmail.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "agustfricke@some.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.UpdateActiveStatusAllUsers("false")
	if err != nil {
		t.Errorf("Expected error to be nil but got: %v", err.Error())
	}
	user, err := database.GetUserByEmail("agustfricke@proton.me")
	if err != nil {
		t.Errorf("Expected error to be nil but got: %v", err.Error())
	}
	if user.IsActive != true {
		t.Errorf("Expected user.IsActive to be false but got: %v", user.IsActive)
	}
	user1, err := database.GetUserByEmail("agustfricke@gmail.me")
	if err != nil {
		t.Errorf("Expected error to be nil but got: %v", err.Error())
	}
	if user1.IsActive != false {
		t.Errorf("Expected user.IsActive to be false but got: %v", user1.IsActive)
	}
	user2, err := database.GetUserByEmail("agustfricke@some.me")
	if err != nil {
		t.Errorf("Expected error to be nil but got: %v", err.Error())
	}
	if user2.IsActive != false {
		t.Errorf("Expected user.IsActive to be false but got: %v", user2.IsActive)
	}

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestUpdateActiveStatus(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
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
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("user not found", func(t *testing.T) {
		err = database.UpdateActiveStatus("99999")
		if err.Error() != "No account found with the id 99999" {
			t.Errorf("Expected error to be 'No account found with the id 99999' but got %v", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		err = database.UpdateActiveStatus("1")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err := database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.IsActive != false {
			t.Errorf("expected IsActive to be false but got %v", user.IsActive)
		}

		err = database.UpdateActiveStatus("1")
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err.Error())
		}
		user, err = database.GetUserByID("1")
		if err != nil {
			t.Errorf("test failed because of GetUserByID(): %v", err.Error())
		}
		if user.IsActive != true {
			t.Errorf("expected IsActive to be true but got %v", user.IsActive)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestGetAdminUsersCountAll(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:    "agustfricke@proton.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "pepe@pepe.me",
		Name:     "Pepe",
		Surname:  "Salamanca",
		Password: "some-password",
		Pc:       "pepe@machine",
		Os:       "Mac",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	t.Run("success", func(t *testing.T) {
		count, err := database.GetAdminUsersCountAll()
		if err != nil {
			t.Errorf("Expected error to be nil but got: %v", err.Error())
			return
		}
		if count != 2 {
			t.Errorf("Expected count to be 2 but got: %v", count)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestGetAdminUsersSearchCount(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:    "agustfricke@proton.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "pepe@pepe.me",
		Name:     "Pepe",
		Surname:  "Salamanca",
		Password: "some-password",
		Pc:       "pepe@machine",
		Os:       "Mac",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}

	t.Run("search", func(t *testing.T) {
		notFound, err := database.GetAdminUsersSearchCount("%not-found-no-no-no%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if notFound != 0 {
			t.Errorf("expected 0 users but got: %v", notFound)
		}
		usersByEmail, err := database.GetAdminUsersSearchCount("%@pepe%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if usersByEmail > 1 {
			t.Errorf("expected 1 or more users but got: %v", usersByEmail)
		}
		usersByName, err := database.GetAdminUsersSearchCount("%Agu%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if usersByName > 1 {
			t.Errorf("expected 1 or more users but got: %v", usersByName)
		}
		usersBySurname, err := database.GetAdminUsersSearchCount("%Sala%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if usersBySurname > 1 {
			t.Errorf("expected 1 or more users but got: %v", usersBySurname)
		}
		usersByCreatedAt, err := database.GetAdminUsersSearchCount("%20%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if usersByCreatedAt != 2 {
			t.Errorf("expected 1 or more users but got: %v", usersByCreatedAt)
		}
		usersByOs, err := database.GetAdminUsersSearchCount("%Mac%", "", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if usersByOs > 1 {
			t.Errorf("expected 1 or more users but got: %v", usersByOs)
		}
	})

	t.Run("is active", func(t *testing.T) {
		isActive, err := database.GetAdminUsersSearchCount("%%", "1", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if isActive != 2 {
			t.Errorf("expected 2 users but got: %v", isActive)
		}
		isActive, err = database.GetAdminUsersSearchCount("%%", "0", "", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if isActive != 0 {
			t.Errorf("expected 0 users but got: %v", isActive)
		}
	})

	t.Run("is admin", func(t *testing.T) {
		isAdmin, err := database.GetAdminUsersSearchCount("%%", "", "0", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if isAdmin != 2 {
			t.Errorf("expected 2 users but got: %v", isAdmin)
		}
		database.DB.Exec(`UPDATE users SET is_admin = 1 WHERE email = 'agustfricke@proton.me';`)
		isAdmin, err = database.GetAdminUsersSearchCount("%%", "", "1", "", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if isAdmin != 1 {
			t.Errorf("expected 1 user but got: %v", isAdmin)
		}
	})

	t.Run("special apps", func(t *testing.T) {
		specialApps, err := database.GetAdminUsersSearchCount("%%", "", "", "0", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if specialApps != 2 {
			t.Errorf("expected 2 users but got: %v", specialApps)
		}

		database.DB.Exec(`UPDATE users SET special_apps = 1 WHERE email = 'agustfricke@proton.me';`)
		specialApps, err = database.GetAdminUsersSearchCount("%%", "", "", "1", "", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if specialApps != 1 {
			t.Errorf("expected 1 user but got: %v", specialApps)
		}
	})

	t.Run("verified", func(t *testing.T) {
		verified, err := database.GetAdminUsersSearchCount("%%", "", "", "", "0", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if verified != 2 {
			t.Errorf("expected 1 user but got: %v", verified)
		}

		database.DB.Exec(`UPDATE users SET verified = 1 WHERE email = 'agustfricke@proton.me';`)
		verified, err = database.GetAdminUsersSearchCount("%%", "", "", "", "1", "", "")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if verified != 1 {
			t.Errorf("expected 1 user but got: %v", verified)
		}
	})

	t.Run("search between dates", func(t *testing.T) {
		users, err := database.GetAdminUsersSearchCount("%%", "", "", "", "", "28/09/2024 01:59:21", "30/09/2050 01:59:21")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if users != 2 {
			t.Errorf("expected 2 users but got: %v", users)
		}

		users, err = database.GetAdminUsersSearchCount("%%", "", "", "", "", "28/09/2004 01:59:21", "28/09/2005 01:59:21")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if users != 0 {
			t.Errorf("expected 0 users but got: %v", users)
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}

func TestGetAdminUsers(t *testing.T) {
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:    "agustfricke@proton.me",
		Name:     "Agust",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	err = database.CreateUser(database.User{
		Email:    "pepe@pepe.me",
		Name:     "Pepe",
		Surname:  "Salamanca",
		Password: "some-password",
		Pc:       "pepe@machine",
		Os:       "Mac",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
		return
	}
	t.Run("limit", func(t *testing.T) {
		users, err := database.GetAdminUsers("%%", "", "", "", "", "", "", 1, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(users) != 1 {
			t.Errorf("expected 1 user but got: %v", len(users))
		}

		users, err = database.GetAdminUsers("%%", "", "", "", "", "", "", 2, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users but got: %v", len(users))
		}
	})

	t.Run("cursor", func(t *testing.T) {
		users, err := database.GetAdminUsers("%%", "", "", "", "", "", "", 1, 1)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(users) != 1 {
			t.Errorf("expected 1 user but got: %v", len(users))
		}
	})

	t.Run("search", func(t *testing.T) {
		notFound, err := database.GetAdminUsers("%not-found-no-no-no%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(notFound) != 0 {
			t.Errorf("expected 0 users but got: %v", len(notFound))
		}
		usersByEmail, err := database.GetAdminUsers("%@pepe%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(usersByEmail) > 1 {
			t.Errorf("expected 1 or more users but got: %v", len(usersByEmail))
		}
		usersByName, err := database.GetAdminUsers("%Agu%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(usersByName) > 1 {
			t.Errorf("expected 1 or more users but got: %v", len(usersByName))
		}
		usersBySurname, err := database.GetAdminUsers("%Sala%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(usersBySurname) > 1 {
			t.Errorf("expected 1 or more users but got: %v", len(usersBySurname))
		}
		usersByCreatedAt, err := database.GetAdminUsers("%20%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(usersByCreatedAt) != 2 {
			t.Errorf("expected 1 or more users but got: %v", len(usersByCreatedAt))
		}
		usersByOs, err := database.GetAdminUsers("%Mac%", "", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(usersByOs) > 1 {
			t.Errorf("expected 1 or more users but got: %v", len(usersByOs))
		}
	})

	t.Run("is active", func(t *testing.T) {
		isActive, err := database.GetAdminUsers("%%", "1", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(isActive) != 2 {
			t.Errorf("expected 2 users but got: %v", len(isActive))
		}
		isActive, err = database.GetAdminUsers("%%", "0", "", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(isActive) != 0 {
			t.Errorf("expected 0 users but got: %v", len(isActive))
		}
	})

	t.Run("is admin", func(t *testing.T) {
		isAdmin, err := database.GetAdminUsers("%%", "", "0", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(isAdmin) != 2 {
			t.Errorf("expected 2 users but got: %v", len(isAdmin))
		}
		database.DB.Exec(`UPDATE users SET is_admin = 1 WHERE email = 'agustfricke@proton.me';`)
		isAdmin, err = database.GetAdminUsers("%%", "", "1", "", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(isAdmin) != 1 {
			t.Errorf("expected 1 user but got: %v", len(isAdmin))
		}
	})

	t.Run("special apps", func(t *testing.T) {
		specialApps, err := database.GetAdminUsers("%%", "", "", "0", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(specialApps) != 2 {
			t.Errorf("expected 2 users but got: %v", len(specialApps))
		}

		database.DB.Exec(`UPDATE users SET special_apps = 1 WHERE email = 'agustfricke@proton.me';`)
		specialApps, err = database.GetAdminUsers("%%", "", "", "1", "", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(specialApps) != 1 {
			t.Errorf("expected 1 user but got: %v", len(specialApps))
		}
	})

	t.Run("verified", func(t *testing.T) {
		verified, err := database.GetAdminUsers("%%", "", "", "", "0", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(verified) != 2 {
			t.Errorf("expected 1 user but got: %v", len(verified))
		}

		database.DB.Exec(`UPDATE users SET verified = 1 WHERE email = 'agustfricke@proton.me';`)
		verified, err = database.GetAdminUsers("%%", "", "", "", "1", "", "", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(verified) != 1 {
			t.Errorf("expected 1 user but got: %v", len(verified))
		}
	})

	t.Run("search between dates", func(t *testing.T) {
		users, err := database.GetAdminUsers("%%", "", "", "", "", "28/09/2024 01:59:21", "30/09/2025 01:59:21", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users but got: %v", len(users))
		}

		users, err = database.GetAdminUsers("%%", "", "", "", "", "01/12/2012%00:00:00", "28/09/2013%00:00:00", 50, 0)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users but got: %v", len(users))
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}

func TestUpdatePassword(t *testing.T) {
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
		err := database.UpdatePassword("new-password", "agustfricke@proton.me")
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
		user, err := database.GetUserByEmail("agustfricke@proton.me")
		if err != nil {
			t.Errorf("test failed because of GetUserByEmail(): %v", err)
		}
		if user.Password == "new-password" {
			t.Errorf("Expected password to be hashed but got 'new-password'")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		err := database.UpdatePassword("new-password", "idont@exist.com")
		if err == nil {
			t.Errorf("expected error but god nil")
		}
		if err.Error() != "No account found with the email idont@exist.com." {
			t.Errorf("expected 'No account found with the email idont@exist.com.' but got %v", err.Error())
		}
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

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
