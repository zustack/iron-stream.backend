package inputs

import (
	"fmt"
	"iron-stream/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginInput(input database.User) (database.User, error) {
	if input.Email == "" {
		return database.User{}, fmt.Errorf("The email is required.")
	}
	if len(input.Email) > 55 {
		return database.User{}, fmt.Errorf("The email should not have more than 55 characters.")
	}
	if input.Password == "" {
		return database.User{}, fmt.Errorf("The password is required.")
	}
	if len(input.Password) > 55 {
		return database.User{}, fmt.Errorf("The password should not have more than 55 characters.")
	}
	if input.Pc == "" {
		return database.User{}, fmt.Errorf("The unique identifier is required. Please ensure that your system's configuration is correct.")
	}
	if len(input.Pc) > 255 {
		return database.User{}, fmt.Errorf("The unique identifier should not have more than 55 characters. Please ensure that your system's configuration is correct.")
	}
	return database.User{
		Email:    input.Email,
		Password: input.Password,
		Pc:       input.Pc,
	}, nil
}

func CleanRegisterInput(input database.User) (database.User, error) {
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
		return database.User{}, fmt.Errorf("Foo No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	if len(input.Os) > 20 {
		return database.User{}, fmt.Errorf("No se pudo crear la cuenta debido a una incompatibilidad con tu sistema operativo.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, fmt.Errorf("Ocurrio un error al generar encriptar la contraseña.")
	}

	return database.User{
		Password: string(hashedPassword),
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		IsAdmin:  false,
		Pc:       input.Pc,
		Os:       input.Os,
	}, nil
}
