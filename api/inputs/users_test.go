package inputs_test

import (
	"iron-stream/api/inputs"
	"strings"
	"testing"
)

func TestSignup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err != nil {
			t.Errorf("test failed because: %v", err)
		}
	})

	t.Run("missing email", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The email is required." {
			t.Errorf("Expected error to be 'The email is required.' but got: %v", err.Error())
		}
	})

	t.Run("email to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    strings.Repeat("a", 56),
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The email should not have more than 55 characters." {
			t.Errorf("Expected error to be 'The email should not have more than 55 characters.' but got: %v", err.Error())
		}
	})

	t.Run("missing password", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The password is required." {
			t.Errorf("Expected error to be 'The password is required.' but got: %v", err.Error())
		}
	})

	t.Run("password to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: strings.Repeat("a", 56),
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The password should not have more than 55 characters." {
			t.Errorf("Expected error to be 'The password should not have more than 55 characters.' but got: %v", err.Error())
		}
	})

	t.Run("password to short", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "123",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The password should have at least 8 characters." {
			t.Errorf("Expected error to be 'The password should have at least 8 characters.' but got: %v", err.Error())
		}
	})

	t.Run("missing name", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The name is required." {
			t.Errorf("Expected error to be 'The name is required.' but got: %v", err.Error())
		}
	})

	t.Run("name to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     strings.Repeat("a", 56),
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The name should not have more than 55 characters." {
			t.Errorf("Expected error to be 'The name should not have more than 55 characters.' but got: %v", err.Error())
		}
	})

	t.Run("missing surname", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The surname is required." {
			t.Errorf("Expected error to be 'The surname is required.' but got: %v", err.Error())
		}
	})

	t.Run("surname to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  strings.Repeat("a", 56),
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The surname should not have more than 55 characters." {
			t.Errorf("Expected error to be 'The surname should not have more than 55 characters.' but got: %v", err.Error())
		}
	})

	t.Run("missing pc", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "",
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The unique identifier is required. Please ensure that your system's configuration is correct." {
			t.Errorf("Expected error to be 'The unique identifier is required. Please ensure that your system's configuration is correct.' but got: %v", err.Error())
		}
	})

	t.Run("pc to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       strings.Repeat("a", 256),
			Os:       "Linux",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct." {
			t.Errorf("Expected error to be 'The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct.'s configuration is correct.' but got: %v", err.Error())
		}
	})

	t.Run("missing os", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       "",
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The os identifier is required. Please ensure that your system's configuration is correct." {
			t.Errorf("Expected error to be 'The os identifier is required. Please ensure that your system's configuration is correct.'s configuration is correct.' but got: %v", err.Error())
		}
	})

	t.Run("pc to long", func(t *testing.T) {
		payload := inputs.SignupInput{
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			Password: "securepassword",
			Pc:       "uniqueID",
			Os:       strings.Repeat("a", 21),
		}
		inputs.Signup(payload)
		_, err := inputs.Signup(payload)
		if err == nil {
			t.Errorf("expected error got nil")
		}
		if err.Error() != "The os identifier should not have more than 20 characters. Please ensure that your system's configuration is correct." {
			t.Errorf("Expected error to be 'The os identifier should not have more than 20 characters. Please ensure that your system's configuration is correct.'s configuration is correct.'s configuration is correct.' but got: %v", err.Error())
		}
	})
}

/*
func TestLoginInput(t *testing.T) {
	t.Skip("Skipping TestLoginInput()")
	tests := []struct {
		name    string
		input   database.User
		want    database.User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid input",
			input:   database.User{Email: "test@example.com", Password: "securepassword", Pc: "uniqueID"},
			want:    database.User{Email: "test@example.com", Password: "securepassword", Pc: "uniqueID"},
			wantErr: false,
		},
		{
			name:    "Empty email",
			input:   database.User{Email: "", Password: "securepassword", Pc: "uniqueID"},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The email is required.",
		},
		{
			name:    "Email too long",
			input:   database.User{Email: strings.Repeat("a", 56), Password: "securepassword", Pc: "uniqueID"},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The email should not have more than 55 characters.",
		},
		{
			name:    "Empty password",
			input:   database.User{Email: "test@example.com", Password: "", Pc: "uniqueID"},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The password is required.",
		},
		{
			name:    "Password too long",
			input:   database.User{Email: "test@example.com", Password: strings.Repeat("a", 56), Pc: "uniqueID"},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The password should not have more than 55 characters.",
		},
		{
			name:    "Empty unique identifier",
			input:   database.User{Email: "test@example.com", Password: "securepassword", Pc: ""},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The unique identifier is required. Please ensure that your system's configuration is correct.",
		},
		{
			name:    "Unique identifier too long",
			input:   database.User{Email: "test@example.com", Password: "securepassword", Pc: strings.Repeat("a", 256)},
			want:    database.User{},
			wantErr: true,
			errMsg:  "The unique identifier should not have more than 255 characters. Please ensure that your system's configuration is correct.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inputs.LoginInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("LoginInput() error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("LoginInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
