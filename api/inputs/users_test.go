package inputs_test

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"strings"
	"testing"
)

func TestLoginInput(t *testing.T) {
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
