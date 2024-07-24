package inputs

import (
	"fmt"
	"iron-stream/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func CleanRegisterInput(input database.User) (database.User, error) {
	if input.Username == "" {
		return database.User{}, fmt.Errorf("El nombre de usuario es requerido.")
	}
	if len(input.Username) > 55 {
		return database.User{}, fmt.Errorf("El nombre de usuario no debe tener más de 55 caracteres.")
	}

	if input.Password == "" {
		return database.User{}, fmt.Errorf("La contraseña es requerida.")
	}
	if len(input.Password) > 55 {
		return database.User{}, fmt.Errorf("La contraseña no debe tener más de 55 caracteres.")
	}

	if input.Email == "" {
		return database.User{}, fmt.Errorf("El email es requerido.")
	}
	if len(input.Email) > 55 {
		return database.User{}, fmt.Errorf("El email no debe tener más de 55 caracteres.")
	}

	if input.Name == "" {
		return database.User{}, fmt.Errorf("El nombre es requerido.")
	}
	if len(input.Name) > 55 {
		return database.User{}, fmt.Errorf("El nombre no debe tener más de 55 caracteres.")
	}

	if input.Surname == "" {
		return database.User{}, fmt.Errorf("El apellido es requerido.")
	}
	if len(input.Surname) > 55 {
		return database.User{}, fmt.Errorf("El apellido no debe tener más de 55 caracteres.")
	}

	if input.Pc == "" {
		return database.User{}, fmt.Errorf("No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	if len(input.Pc) > 255 {
		return database.User{}, fmt.Errorf("No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	if input.Os == "" {
		return database.User{}, fmt.Errorf("No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	if len(input.Os) > 20 {
		return database.User{}, fmt.Errorf("No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	is_admin := false
	pc := input.Pc
	os := input.Os
	if input.Username == "admin" {
		is_admin = true
		pc = ""
		os = ""
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, fmt.Errorf("Error hashing password")
	}

	/*
		assert.True(t, user.ValidatePassword(pw))
		assert.False(t, user.ValidatePassword("hunter2005"))
	*/

	return database.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		IsAdmin:  is_admin,
		Pc:       pc,
		Os:       os,
	}, nil
}
