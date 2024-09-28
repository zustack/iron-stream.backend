package inputs

import (
	"fmt"
	"iron-stream/internal/database"
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
	if input.Pc == "" {
		return database.User{}, fmt.Errorf("The unique identifier is required. Please ensure that your system's configuration is correct.")
	}
	if len(input.Pc) > 255 {
		return database.User{}, fmt.Errorf("The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct.")
	}
	return database.User{
		Email:    input.Email,
		Password: input.Password,
		Pc:       input.Pc,
	}, nil
}

type SignupInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
	Pc       string `json:"pc"`
	Os       string `json:"os"`
}

func Signup(input SignupInput) (SignupInput, error) {
	if input.Email == "" {
		return SignupInput{}, fmt.Errorf("The email is required.")
	}
	if len(input.Email) > 55 {
		return SignupInput{}, fmt.Errorf("The email should not have more than 55 characters.")
	}

	if input.Password == "" {
		return SignupInput{}, fmt.Errorf("The password is required.")
	}
	if len(input.Password) > 55 {
		return SignupInput{}, fmt.Errorf("The password should not have more than 55 characters.")
	}
	if len(input.Password) < 8 {
		return SignupInput{}, fmt.Errorf("The password should have at least 8 characters.")
	}

	if input.Name == "" {
		return SignupInput{}, fmt.Errorf("The name is required.")
	}
	if len(input.Name) > 55 {
		return SignupInput{}, fmt.Errorf("The name should not have more than 55 characters.")
	}

	if input.Surname == "" {
		return SignupInput{}, fmt.Errorf("The surname is required.")
	}
	if len(input.Surname) > 55 {
		return SignupInput{}, fmt.Errorf("The surname should not have more than 55 characters.")
	}

	if input.Pc == "" {
		return SignupInput{}, fmt.Errorf("The unique identifier is required. Please ensure that your system's configuration is correct.")
	}
	if len(input.Pc) > 255 {
		return SignupInput{}, fmt.Errorf("The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct.")
	}

	if input.Os == "" {
		return SignupInput{}, fmt.Errorf("The os identifier is required. Please ensure that your system's configuration is correct.")
	}
	if len(input.Os) > 20 {
		return SignupInput{}, fmt.Errorf("The os identifier should not have more than 20 characters. Please ensure that your system's configuration is correct.")
	}

	return SignupInput{
		Password: input.Password,
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Pc:       input.Pc,
		Os:       input.Os,
	}, nil
}
