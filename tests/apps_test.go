package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"iron-stream/api"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func generateToken(t *testing.T, user_id int64, is_admin bool) string {
	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := time.Hour * 24 * 30
	claims["sub"] = user_id
	claims["is_admin"] = is_admin
	claims["exp"] = now.Add(expDuration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(utils.GetEnv("SECRET_KEY")))
	if err != nil {
		t.Fatal(err)
	}
	return tokenString
}

func TestCreateApp(t *testing.T) {
	app := api.Setup()
	err := database.ExecuteSQLFile("../test_sqlite.db", "../tables.sql")
	if err != nil {
		t.Fatal(err)
	}
	database.ConnectDB("DB_TEST_PATH")

	id, err := database.CreateUser(database.User{
		Username: "admin",
		Password: "admin-password",
		Email:    "admin@admin.com",
		Name:     "admin",
		Surname:  "admin",
		IsAdmin:  true,
		Pc:       "admin-pc",
		Os:       "admin-os",
	})

	if err != nil {
		t.Fatal(err)
	}

	adminToken := generateToken(t, id, true)

	id, err = database.CreateUser(database.User{
		Username: "normal",
		Password: "normal-password",
		Email:    "normal@normal.com",
		Name:     "normal",
		Surname:  "normal",
		IsAdmin:  false,
		Pc:       "normal-pc",
		Os:       "normal-os",
	})

	if err != nil {
		t.Fatal(err)
	}

	normalUserToken := generateToken(t, id, false)

	tests := []struct {
		description  string
		payload      database.App
		expectedCode int
		expectedBody string
		authHeader   string
	}{
		{
			description: "create app successfully",
			payload: database.App{
				Name:        "test-app",
				ProcessName: "test-process",
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusCreated,
			expectedBody: `{"id":1}`,
		},

		// missing name
		{
			description: "missing name",
			payload: database.App{
				Name:        "",
				ProcessName: "test-process",
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre es requerido."}`,
		},
		// name to long
		{
			description: "name to long",
			payload: database.App{
				Name:        strings.Repeat("a", 56),
				ProcessName: "test-process",
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre no debe tener más de 55 caracteres."}`,
		},
		// missing process name
		{
			description: "missing process name",
			payload: database.App{
				Name:        "test-app",
				ProcessName: "",
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre del proceso es requerido."}`,
		},
		// process name to long
		{
			description: "process name to long",
			payload: database.App{
				Name:        "test-app",
				ProcessName: strings.Repeat("a", 56),
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre del proceso no debe tener más de 55 caracteres."}`,
		},
		// missing os
		{
			description: "missing os",
			payload: database.App{
				Name:        "test-app",
				ProcessName: "test-process",
				Os:          "",
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El sistema operativo es requerido."}`,
		},
		// os to long
		{
			description: "os to long",
			payload: database.App{
				Name:        "test-app",
				ProcessName: "test-process",
				Os:          strings.Repeat("a", 56),
				IsActive:    true,
			},
			authHeader:   "Bearer " + adminToken,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El sistema operativo no debe tener más de 55 caracteres."}`,
		},

		{
			description: "not admin user",
			payload: database.App{
				Name:        "test-app",
				ProcessName: "test-process",
				Os:          "test-os",
				IsActive:    true,
			},
			authHeader:   "Bearer " + normalUserToken,
			expectedCode: fiber.StatusForbidden,
			expectedBody: `{"error":"No eres administrador."}`,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.payload)
			req, _ := http.NewRequest("POST", "/apps/create", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", test.authHeader)

			res, _ := app.Test(req, -1)
			assert.Equal(t, test.expectedCode, res.StatusCode, test.description)

			resBody, _ := io.ReadAll(res.Body)
			assert.JSONEq(t, test.expectedBody, string(resBody), test.description)
		})
	}
}
