package database

import (
	"fmt"
)

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin"`
	IsActive   bool   `json:"is_active"`
	EmailToken int    `json:"email_token"`
	Verified   bool   `json:"verified"`
	Courses    string `json:"courses"`
	Pc         string `json:"pc"`
	Os         string `json:"os"`
	CreatedAt  string `json:"created_at"`
}

func CreateUser(u User) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO users
  (username, email, name, surname, password, is_admin, email_token, pc, os) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Username, u.Email, u.Name, u.Surname, u.Password, u.IsAdmin, u.EmailToken, u.Pc, u.Os)

	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	return id, nil
}
