package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"iron-stream/api"
	"iron-stream/internal/database"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyEmail(t *testing.T) {
	app := api.Setup()
	if app == nil {
		t.Fatal("Failed to initialize app")
	}
	database.ConnectDB("DB_DEV_PATH")

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
