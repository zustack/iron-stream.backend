package handlers_test

import (
	"bytes"
	"encoding/json"
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

	tests := []struct {
		name         string
		payload      string
		expectedCode int
	}{
		{
			name:         "success admin user",
			payload:      `{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "agust@ubuntu"}`,
			expectedCode: 200,
		},
		{
			name:         "json error",
			payload:      `{"email": "", "password": "some-password", "pc": "some-pc}`,
			expectedCode: 400,
		},
		{
			name:         "normal user wrong pc",
			payload:      `{"email": "agustfricke@protonmail.com", "password": "some-password", "pc": "wrong-pc"}`,
			expectedCode: 401,
		},
		{
			name:         "admin user wrong pc",
			payload:      `{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "wrong-pc"}`,
			expectedCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var body io.Reader
			body = bytes.NewBufferString(test.payload)
			req, _ := http.NewRequest("POST", "/users/login", body)
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req)
			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode)

			responseBody, _ := io.ReadAll(res.Body)

			if test.expectedCode == 200 {
				var responseData SuccessLoginResponse
				if err := json.Unmarshal(responseBody, &responseData); err != nil {
					t.Fatalf("Error unmarshaling response: %v", err)
				}
				assert.NotEmpty(t, responseData.Token)
				assert.NotEmpty(t, responseData.UserID)
				assert.NotEmpty(t, responseData.Exp)
				assert.NotEmpty(t, responseData.FullName)
				if test.name == "success admin user" || test.name == "admin user wrong pc" {
					assert.True(t, responseData.IsAdmin)
				} else {
					assert.False(t, responseData.IsAdmin)
				}
			}
		})
	}
}
