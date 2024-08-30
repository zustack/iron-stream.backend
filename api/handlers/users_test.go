package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"iron-stream/api"
	"iron-stream/internal/database"

	"github.com/stretchr/testify/assert"
)

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

	id, err := database.CreateUser(database.User{
		Email:    "agustfricke@gmail.com",
		Name:     "Agustin",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
	})
	if err != nil {
		t.Errorf("test failed because of CreateUser(): %v", err)
	}
	err = database.UpdateAdminStatus(fmt.Sprintf("%d", id), "true")
	if err != nil {
		t.Errorf("test failed because of UpdateAdminStatus(): %v", err)
	}

	_, err = database.CreateUser(database.User{
		Email:    "agustfricke@protonmail.com",
		Name:     "Agustin",
		Surname:  "Fricke",
		Password: "some-password",
		Pc:       "agust@ubuntu",
	})

	t.Run("success admin user", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "agust@ubuntu"}`)
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
		body = bytes.NewBufferString(`{"email": "agustfricke@protonmail.com", "password": "some-password", "pc": "agust@ubuntu"}`)
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

	t.Run("success admin user with wrong pc", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "wrong-pc"}`)
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

	t.Run("error normal user with wrong pc", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`{"email": "agustfricke@protonmail.com", "password": "some-password", "pc": "wrong-pc"}`)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 401, res.StatusCode)
	})

	t.Run("error json", func(t *testing.T) {
		var body io.Reader
		body = bytes.NewBufferString(`{"email: "agustfricke@gmail.com", "password": 1, "pc": "agust@ubuntu"}`)
		req, _ := http.NewRequest("POST", "/users/login", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})

	_, err = database.DB.Exec(`DELETE FROM users WHERE email IN ('agustfricke@gmail.com', 'agustfricke@protonmail.com')`)
	if err != nil {
		t.Fatalf("failed to teardown test database: %v", err)
	}

}
