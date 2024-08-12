package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
	"strings"
)

type User struct {
	ID         int64  `json:"id"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	IsAdmin    bool   `json:"is_admin"`
  SpecialApps bool   `json:"special_apps"`
	IsActive   bool   `json:"is_active"`
	EmailToken int    `json:"email_token"`
	Verified   bool   `json:"verified"`
	Courses    string `json:"courses"`
	Pc         string `json:"pc"`
	Os         string `json:"os"`
	CreatedAt  string `json:"created_at"`
}

func UpdateEmailToken(email string, email_token int) error {
  _, err := DB.Exec(`UPDATE users SET email_token = ? WHERE email = ?`, email_token, email)
  if err != nil {
    return fmt.Errorf("UpdateEmailToken: %v", err)
  }
  return nil
}

func UpdatePassword(password, email string) error {
    // Trim any whitespace from email
    email = strings.TrimSpace(email)

    // Execute the update
    result, err := DB.Exec(`UPDATE users SET password = ? WHERE LOWER(email) = LOWER(?)`, password, email)
    if err != nil {
        return fmt.Errorf("UpdatePassword: %v", err)
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("UpdatePassword: couldn't get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("UpdatePassword: no user found with email %s", email)
    }

    return nil
}

func DeleteAccount(email string) error {
  _, err := DB.Exec(`DELETE FROM users WHERE email = ?`, email)
  if err != nil {
    return fmt.Errorf("DeleteAccount: %v", err)
  }
  return nil
}

func VerifyAccount(userID int64) error {
  _, err := DB.Exec(`UPDATE users SET verified = true WHERE id = ?`, userID)
  if err != nil {
    return fmt.Errorf("EditSortCourses: %v", err)
  }
  return nil
}

func GetUserByID(id string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
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
	if err := row.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
		&u.Courses, &u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("GetUserByEmail: %s: no such user", email)
		}
		return u, fmt.Errorf("GetUserByEmail: %s: %v", email, err)
	}
	return u, nil
}

func CreateUser(u User) (int64, error) {
  date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO users
  (email, name, surname, password, is_admin, email_token, pc, os, created_at) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email, u.Name, u.Surname, u.Password, true, u.EmailToken, u.Pc, u.Os, date)

	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	return id, nil
}
