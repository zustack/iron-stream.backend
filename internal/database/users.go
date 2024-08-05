package database

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	IsAdmin    bool   `json:"is_admin"`
	IsActive   bool   `json:"is_active"`
	EmailToken int    `json:"email_token"`
	Verified   bool   `json:"verified"`
	Courses    string `json:"courses"`
	Pc         string `json:"pc"`
	Os         string `json:"os"`
	CreatedAt  string `json:"created_at"`
}

func GetUserByID(id string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.IsActive, &u.EmailToken, &u.Verified,
		&u.Courses, &u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("GetUserByID%s: no such user", id)
		}
		return u, fmt.Errorf("GetUserByID: %s: %v", id, err)
	}
	return u, nil
}

func GetUserByEmail(email string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT * FROM users WHERE email = ?`, email)
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.IsActive, &u.EmailToken, &u.Verified,
		&u.Courses, &u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
      return u, fmt.Errorf("GetUserByEmail: %s: no such user", email)
		}
		return u, fmt.Errorf("GetUserByEmail: %s: %v", email, err)
	}
	return u, nil
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
