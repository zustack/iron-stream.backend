package inputs

import (
	"fmt"
)

type VerifyEmailInput struct {
	Email      string `json:"email"`
	EmailToken int `json:"email_token"`
}

func VerifyEmail(input VerifyEmailInput) (VerifyEmailInput, error) {
	if input.Email == "" {
		return VerifyEmailInput{}, fmt.Errorf("The email is required.")
	}
	if input.EmailToken == 0 {
		return VerifyEmailInput{}, fmt.Errorf("The email token is required.")
	}
	if input.EmailToken < 100000  {
		return VerifyEmailInput{}, fmt.Errorf("The email token is to small.")
	}
	if input.EmailToken > 999999 {
		return VerifyEmailInput{}, fmt.Errorf("The email token is to large.")
	}
	return VerifyEmailInput{
		Email:      input.Email,
		EmailToken: input.EmailToken,
	}, nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
	Pc       string `json:"pc"`
	Os       string `json:"os"`
}

func Login(input LoginInput) (LoginInput, error) {
	if input.Email == "" {
		return LoginInput{}, fmt.Errorf("The email is required.")
	}
	if len(input.Email) > 55 {
		return LoginInput{}, fmt.Errorf("The email should not have more than 55 characters.")
	}
	if input.Password == "" {
		return LoginInput{}, fmt.Errorf("The password is required.")
	}
	if input.Pc == "" {
		return LoginInput{}, fmt.Errorf("The unique identifier is required. Please ensure that your system's configuration is correct.")
	}
	if len(input.Pc) > 255 {
		return LoginInput{}, fmt.Errorf("The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct.")
	}
	return LoginInput{
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
