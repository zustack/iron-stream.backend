package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iron-stream/api"
	"iron-stream/api/handlers"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserStatistics(t *testing.T) {
	t.Skip("Skip")
}

func TestGetCurrentUser(t *testing.T) {
	t.Skip("Skip")
}

func TestUpdateAdminStatus(t *testing.T) {
	t.Skip("Skip")
}

func TestUpdateSpecialAppUser(t *testing.T) {
	t.Skip("Skip")
}

func TestUpdateActiveStatusAllUsers(t *testing.T) {
	t.Skip("Skip")
}

func TestUpdateActiveStatus(t *testing.T) {
	t.Skip("Skip")
}

func TestAdminUsers(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
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
		Email:      "agustfricke@gmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	database.DB.Exec(`UPDATE users SET is_admin = true WHERE email = 'agustfricke@gmail.com';`)
	token, _, err := utils.MakeJWT(1)
	if err != nil {
		t.Error("Test failed because of utils.MakeJWT")
	}
	err = database.CreateUser(database.User{
		Email:      "test@test.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	unToken, _, err := utils.MakeJWT(2)
	if err != nil {
		t.Error("Test failed because of utils.MakeJWT")
	}

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/admin", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		res, _ := app.Test(req, -1)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		var response handlers.AdminUsersResponse
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		if len(response.Data) != 2 {
			t.Errorf("expected 2 users but got %v", len(response.Data))
		}
		if response.TotalCount != 2 {
			t.Errorf("expected 2 total count but got %v", response.TotalCount)
		}
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("success with cursor", func(t *testing.T) {
		for i := 1; i <= 150; i++ {
			email := "user" + strconv.Itoa(i) + "@test.com"
			user := database.User{
				Email:      email,
				Name:       "Agustin",
				Surname:    "Fricke",
				Password:   "some-password",
				Pc:         "agust@ubuntu",
				Os:         "Linux",
				EmailToken: 123456,
			}
			err := database.CreateUser(user)
			if err != nil {
				t.Errorf("test failed because of CreateUser(): %v", err)
			}
		}

		req, _ := http.NewRequest("GET", "/users/admin?cursor=53&q=&a=&admin=&special=&verified=&from=&to=", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		res, _ := app.Test(req, -1)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		var response handlers.AdminUsersResponse
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}

		if len(response.Data) != 50 {
			t.Errorf("expected 50 users but got %v", len(response.Data))
		}
		if response.TotalCount != 152 {
			t.Errorf("expected 152 total count but got %v", response.TotalCount)
		}
		if *response.NextID != 103 {
			t.Errorf("expected NextID to be nil but got %v", response.NextID)
		}
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/admin", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+unToken)

		res, _ := app.Test(req, -1)
		defer res.Body.Close()

		assert.Equal(t, 403, res.StatusCode)
	})

}

func TestUpdatePassword(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
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
		Email:      "agustfricke@gmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	database.DB.Exec(`UPDATE users SET is_admin = true WHERE email = 'agustfricke@gmail.com';`)
	token, _, err := utils.MakeJWT(1)
	if err != nil {
		t.Error("Test failed because of utils.MakeJWT")
	}

	err = database.CreateUser(database.User{
		Email:      "agustfricke@protonmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	t.Run("success", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "password": "123456"
    }
    `)
		req, _ := http.NewRequest("PUT", "/users/update/password", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("bad request", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com"
      "password": "123456"
    }
    `)
		req, _ := http.NewRequest("PUT", "/users/update/password", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})

	t.Run("forbidden", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@protonmail.com", 
      "password": "123456"
    }
    `)
		req, _ := http.NewRequest("PUT", "/users/update/password", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 403, res.StatusCode)
	})
}

func TestDeleteAccountByEmail(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
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
		Email:      "agustfricke@gmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	database.DB.Exec(`UPDATE users SET is_admin = true WHERE email = 'agustfricke@gmail.com';`)
	adminToken, _, err := utils.MakeJWT(1)
	if err != nil {
		t.Error("Test failed because of utils.MakeJWT")
	}

	err = database.CreateUser(database.User{
		Email:      "agustfricke@protonmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	token, _, err := utils.MakeJWT(2)
	if err != nil {
		t.Error("Test failed because of utils.MakeJWT")
	}

	err = database.CreateUser(database.User{
		Email:      "test@test.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	t.Run("success normal user", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/delete/account/by/email/agustfricke@protonmail.com", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("success admin", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/delete/account/by/email/test@test.com", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("normal user try to delete other user account", func(t *testing.T) {
		err = database.CreateUser(database.User{
			Email:      "other@gmail.com",
			Name:       "Agustin",
			Surname:    "Fricke",
			Password:   "some-password",
			Pc:         "agust@ubuntu",
			Os:         "Linux",
			EmailToken: 123456,
		})
		if err != nil {
			t.Errorf("test failed because of CreateUser(): %v", err)
		}
		token, _, err := utils.MakeJWT(4)
		if err != nil {
			t.Error("Test failed because of utils.MakeJWT")
		}

		err = database.CreateUser(database.User{
			Email:      "other2@gmail.com",
			Name:       "Agustin",
			Surname:    "Fricke",
			Password:   "some-password",
			Pc:         "agust@ubuntu",
			Os:         "Linux",
			EmailToken: 123456,
		})
		if err != nil {
			t.Errorf("test failed because of CreateUser(): %v", err)
		}

		req, _ := http.NewRequest("DELETE", "/users/delete/account/by/email/other2@gmail.com", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 403, res.StatusCode)
	})

	t.Run("missing jwt", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/delete/account/by/email/agustfricke@protonmail.com", nil)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}

func TestResendEmailToken(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:      "agustfricke@gmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	req, _ := http.NewRequest("POST", "/users/resend/email/token/agustfricke@gmail.com", nil)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req, -1)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}

func TestVerifyEmail(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
	database.ConnectDB("DB_DEV_PATH")
	err := database.CreateUser(database.User{
		Email:      "agustfricke@gmail.com",
		Name:       "Agustin",
		Surname:    "Fricke",
		Password:   "some-password",
		Pc:         "agust@ubuntu",
		Os:         "Linux",
		EmailToken: 123456,
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	t.Run("success", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "email_token": 123456
    }
    `)
		req, _ := http.NewRequest("POST", "/users/verify", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)

		var responseData SuccessLoginResponse
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		assert.NotEmpty(t, responseData.Token)
		assert.NotEmpty(t, responseData.UserID)
		assert.NotEmpty(t, responseData.Exp)
		assert.NotEmpty(t, responseData.FullName)
		assert.False(t, responseData.IsAdmin)
	})

	t.Run("bad request", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com"
      "email_token": 123456
    }
    `)
		req, _ := http.NewRequest("POST", "/users/verify", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})

	t.Run("token not valid", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "email_token": 923450
    }
    `)
		req, _ := http.NewRequest("POST", "/users/verify", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestSignup(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
	database.ConnectDB("DB_DEV_PATH")

	t.Run("success", func(t *testing.T) {
		body := bytes.NewBufferString(`
      {
        "email": "agustfricke@protonmail.com",
        "name": "Agustin",
        "surname": "Fricke",
        "password": "some-password",
        "pc": "agust@ubuntu",
        "os": "Linux"
      }
    `)

		req, err := http.NewRequest("POST", "/users/signup", body)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		if err != nil {
			t.Errorf("Test failed because: %v", err)
			return
		}
		defer res.Body.Close()

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
			return
		}

		var errorResponse ErrorResponse
		if res.StatusCode != 201 {
			if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
				t.Fatalf("Failed to unmarshal error response: %v", err)
			}
			t.Errorf("expected status code 201 but got %d, error: %s", res.StatusCode, errorResponse.Error)
			return
		}

		_, err = database.DB.Exec(`DELETE FROM users`)
		if err != nil {
			t.Fatalf("failed to teardown test database: %v", err)
		}
	})

	t.Run("bad body", func(t *testing.T) {
		body := bytes.NewBufferString(`
      {
        "email": "agustfricke@protonmail.com"
        "name": "Agustin"
        "surname": "Fricke"
        "password": "some-password"
        "pc": "agust@ubuntu"
        "os": "Linux"
      }
    `)

		req, err := http.NewRequest("POST", "/users/signup", body)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		if err != nil {
			t.Errorf("Test failed because: %v", err)
			return
		}
		defer res.Body.Close()

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
			return
		}

		var errorResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			t.Fatalf("Failed to unmarshal error response: %v", err)
		}

		if res.StatusCode != 400 {
			t.Errorf("expected status code 400 but got %d, %s", res.StatusCode, errorResponse.Error)
			return
		}

	})
}

type SuccessLoginResponse struct {
	Token    string `json:"token"`
	UserID   int64  `json:"userId"`
	IsAdmin  bool   `json:"isAdmin"`
	Exp      int64  `json:"exp"`
	FullName string `json:"fullName"`
}

func TestLogin(t *testing.T) {
	app := api.Setup()
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
		Email:    "agustfricke@gmail.com",
		Name:     "Agustin",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
		Os:       "Linux",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}

	err = database.UpdateAdminStatus("1", "true")
	if err != nil {
		t.Errorf("test failed because of UpdateAdminStatus(): %v", err)
	}

	err = database.CreateUser(database.User{
		Email:    "agustfricke@protonmail.com",
		Name:     "Agustin",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
	})

	t.Run("success admin user", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "password": "some-password", 
      "pc": "agust@ubuntu"
    }
    `)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)

		var responseData SuccessLoginResponse
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		assert.NotEmpty(t, responseData.Token)
		assert.NotEmpty(t, responseData.UserID)
		assert.NotEmpty(t, responseData.Exp)
		assert.NotEmpty(t, responseData.FullName)
		assert.True(t, responseData.IsAdmin)
	})

	t.Run("success normal user", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@protonmail.com", 
      "password": "some-password", 
      "pc": "agust@ubuntu"
    }
    `)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)

		var responseData SuccessLoginResponse
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		assert.NotEmpty(t, responseData.Token)
		assert.NotEmpty(t, responseData.UserID)
		assert.NotEmpty(t, responseData.Exp)
		assert.NotEmpty(t, responseData.FullName)
		assert.False(t, responseData.IsAdmin)
	})

	t.Run("incorrect pc", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "password": "some-password", 
      "pc": "wrong-pc"
    }
    `)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)

		var responseData ErrorResponse
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		if responseData.Error != "Incorrect unique identifier, please try again." {
			t.Errorf("expected error to be 'Incorrect unique identifier, please try again.' but got: %v", responseData.Error)
		}

	})

	t.Run("incorrect password", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email": "agustfricke@gmail.com", 
      "password": "wrong-password", 
      "pc": "agust@ubuntu"
    }
    `)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)

		var responseData ErrorResponse
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		if responseData.Error != "Incorrect password, please try again." {
			t.Errorf("expected error to be 'Incorrect password, please try again.' but got: %v", responseData.Error)
		}

	})

	t.Run("error json", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`
    {
      "email: "agustfricke@gmail.com", 
      "password": 1, 
      "pc": "agust@ubuntu"
    }`)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})

	_, err = database.DB.Exec(`DELETE FROM users;`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}
}
