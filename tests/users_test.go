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
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	app := api.Setup()
	err := database.ExecuteSQLFile("../test_sqlite.db", "../tables.sql")
	if err != nil {
		t.Fatal(err)
	}
	database.ConnectDB("DB_TEST_PATH")
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

func TestLogin(t *testing.T) {
	app := api.Setup()
	err := database.ExecuteSQLFile("../test_sqlite.db", "../tables.sql")
	if err != nil {
		t.Fatal(err)
	}
	database.ConnectDB("DB_TEST_PATH")

	adminUserHashedPassword, err := bcrypt.GenerateFromPassword([]byte("a"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	_, err = database.CreateUser(database.User{
		Username: "admin",
		Password: string(adminUserHashedPassword),
		Email:    "a@a.com",
		Name:     "a",
		Surname:  "a",
		IsAdmin:  false,
    // si el usuario es admin, en el request handler, el pc se anula!
		Pc:       "",
		Os:       "a",
	})

	if err != nil {
		t.Fatal(err)
	}


	normalUserHashedPassword, err := bcrypt.GenerateFromPassword([]byte("n"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	_, err = database.CreateUser(database.User{
		Username: "n",
		Password: string(normalUserHashedPassword),
		Email:    "n@n.com",
		Name:     "n",
		Surname:  "n",
		IsAdmin:  false,
		Pc:       "n",
		Os:       "n",
	})

	if err != nil {
		t.Fatal(err)
	}

  tests := []struct {
		description  string
		payload      database.User
		expectedCode int
		expectToken  bool
		expectedBody string
	}{
		{
			description: "normal user login successfully",
			payload: database.User{
				Username: "n",
				Password: "n",
				Pc:       "n",
			},
			expectedCode: fiber.StatusOK,
			expectToken:  true,
		},

		{
			description: "invalid pc",
			payload: database.User{
				Username: "n",
				Password: "n",
				Pc:       "some-wrong-pc",
			},
			expectedCode: fiber.StatusUnauthorized,
			expectToken:  false,
			expectedBody: `{"error":"Esta cuenta esta registrada en otra computadora."}`,
		},

		{
			description: "username not found",
			payload: database.User{
				Username: "some-wrong-username",
				Password: "n",
				Pc:       "n",
			},
			expectedCode: fiber.StatusUnauthorized,
			expectToken:  false,
			expectedBody: `{"error":"No se encontro el usuario con el nombre de usuario proporcionado."}`,
		},

		{
			description: "incorrect password",
			payload: database.User{
				Username: "n",
				Password: "some-wrong-password",
				Pc:       "n",
			},
			expectedCode: fiber.StatusUnauthorized,
			expectToken:  false,
			expectedBody: `{"error":"La contraseña es incorrecta."}`,
		},

		{
			description: "admin login",
			payload: database.User{
				Username: "admin",
				Password: "a",
				Pc:       "i dont care about pc when is admin",
			},
			expectedCode: fiber.StatusOK,
			expectToken:  true,
		},

		{
			description: "missing username",
			payload: database.User{
				Username: "",
				Password: "a",
				Pc:       "a",
			},
			expectedCode: fiber.StatusBadRequest,
			expectToken:  false,
			expectedBody: `{"error":"El nombre de usuario es requerido."}`,
		},

		{
			description: "missing password",
			payload: database.User{
				Username: "a",
				Password: "",
				Pc:       "a",
			},
			expectedCode: fiber.StatusBadRequest,
			expectToken:  false,
			expectedBody: `{"error":"La contraseña es requerida."}`,
		},

		{
			description: "missing pc normal user",
			payload: database.User{
				Username: "a",
				Password: "a",
				Pc:       "",
			},
			expectedCode: fiber.StatusBadRequest,
			expectToken:  false,
			expectedBody: `{"error":"Ocurrio un error debido a una incompatibilidad con tu sistema operativo."}`,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.payload)
			body := bytes.NewBuffer(jsonBody)
			req, _ := http.NewRequest("POST", "/login", body)
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)

			assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

			resBody, _ := io.ReadAll(res.Body)

			var response map[string]interface{}
			err := json.Unmarshal(resBody, &response)
			assert.NoError(t, err, "El cuerpo de la respuesta debería ser JSON válido")

			if test.expectToken {
				token, exists := response["token"]
				assert.True(t, exists, "The response should contain a 'token' field")
				assert.NotEmpty(t, token, "The token should not be empty")
			} else {
				_, exists := response["token"]
				assert.False(t, exists, "The response should not contain a 'token' field")
				assert.JSONEqf(t, test.expectedBody, string(resBody), test.description)
			}
		})
	}
}
