package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"iron-stream/api"
	"iron-stream/internal/database"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	tests := []struct {
		description  string
		payload      database.User
		expectedCode int
		expectedBody string
	}{
		{
			description: "create a normal user successfully",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusCreated,
			expectedBody: `{"id":1}`,
		},

		{
			description: "create admin user successfully",
			payload: database.User{
				Username: "admin",
				Password: "some-password",
				Email:    "admin@email.com",
				Name:     "admin-name",
				Surname:  "admin-surname",
				Pc:       "admin-pc",
				Os:       "admin-os",
			},
			expectedCode: fiber.StatusCreated,
			expectedBody: `{"id":2, "pc": "", "os": "", "is_admin": true}`,
		},

		{
			description: "missing username",
			payload: database.User{
				Username: "",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre de usuario es requerido."}`,
		},
		{
			description: "username to long",
			payload: database.User{
				Username: strings.Repeat("a", 56),
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre de usuario no debe tener más de 55 caracteres."}`,
		},

		{
			description: "missing password",
			payload: database.User{
				Username: "some-username",
				Password: "",
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"La contraseña es requerida."}`,
		},
		{
			description: "password to long",
			payload: database.User{
				Username: "some-username",
				Password: strings.Repeat("a", 56),
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"La contraseña no debe tener más de 55 caracteres."}`,
		},

		{
			description: "missing email",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El email es requerido."}`,
		},
		{
			description: "email to long",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    strings.Repeat("a", 56),
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El email no debe tener más de 55 caracteres."}`,
		},

		{
			description: "missing name",
			payload: database.User{
				Username: "test-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre es requerido."}`,
		},
		{
			description: "name to long",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     strings.Repeat("a", 56),
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El nombre no debe tener más de 55 caracteres."}`,
		},

		{
			description: "missing surname",
			payload: database.User{
				Username: "test-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  "",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El apellido es requerido."}`,
		},
		{
			description: "surname to long",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  strings.Repeat("a", 56),
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"El apellido no debe tener más de 55 caracteres."}`,
		},

		{
			description: "missing pc",
			payload: database.User{
				Username: "test-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  "some-surname",
				Pc:       "",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo."}`,
		},
		{
			description: "pc to long",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  "some-surname",
				Pc:       strings.Repeat("a", 256),
				Os:       "test-os",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo."}`,
		},

		{
			description: "missing os",
			payload: database.User{
				Username: "test-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  "some-surname",
				Pc:       "test-pc",
				Os:       "",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo."}`,
		},
		{
			description: "os to long",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "some-name",
				Surname:  "some-surname",
				Pc:       "test-pc",
				Os:       strings.Repeat("a", 21),
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":"No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo."}`,
		},

		{
			description: "email already taken",
			payload: database.User{
				Username: "new-username",
				Password: "some-password",
				Email:    "some@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusInternalServerError,
			expectedBody: `{"error":"El correo electrónico esta tomado."}`,
		},
		{
			description: "username already taken",
			payload: database.User{
				Username: "some-username",
				Password: "some-password",
				Email:    "new@email.com",
				Name:     "test-name",
				Surname:  "test-surname",
				Pc:       "test-pc",
				Os:       "test-os",
			},
			expectedCode: fiber.StatusInternalServerError,
			expectedBody: `{"error":"El nombre de usuario esta tomado."}`,
		},
	}

	app := api.Setup()
	database.ConnectDB("DB_TEST_PATH")

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			var body io.Reader
			jsonBody, _ := json.Marshal(test.payload)
			body = bytes.NewBuffer(jsonBody)
			req, _ := http.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)

			assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

			resBody, _ := io.ReadAll(res.Body)
			assert.JSONEqf(t, test.expectedBody, string(resBody), test.description)
		})
	}
}
