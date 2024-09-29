package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int64  `json:"id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	IsAdmin     bool   `json:"is_admin"`
	SpecialApps bool   `json:"special_apps"`
	IsActive    bool   `json:"is_active"`
	EmailToken  int    `json:"email_token"`
	Verified    bool   `json:"verified"`
	Pc          string `json:"pc"`
	Os          string `json:"os"`
	CreatedAt   string `json:"created_at"`
}

type UserStatisticsResponse struct {
	CreatedAt string `json:"created_at"`
	Linux     int64  `json:"linux"`
	Windows   int64  `json:"windows"`
	Mac       int64  `json:"mac"`
	Total     int64  `json:"total"`
}

func GetUserCount(day, os string) (int, error) {
	var count int
	day = "%" + day + "%"
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE created_at LIKE ? AND os = ?;`, day, os).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return count, nil
}

func UpdateEmailToken(email string, email_token int) error {
	result, err := DB.Exec(`UPDATE users SET email_token = ? WHERE email = ?`, email_token, email)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found with the email %s", email)
	}
	return nil
}

func GetUserByID(id string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
		&u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("No user found with id %s", id)
		}
		return u, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return u, nil
}

func UpdateAdminStatus(userId, isAdmin string) error {
	result, err := DB.Exec(`UPDATE users SET is_admin = ? WHERE id = ?`, isAdmin, userId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found with the id %s", userId)
	}
	return nil
}

func UpdateUserSpecialApps(userId, special_apps string) error {
	result, err := DB.Exec(`UPDATE users SET special_apps = ? WHERE id = ?`, special_apps, userId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found with the id %s", userId)
	}
	return nil
}

func UpdateActiveStatusAllUsers(isActive string) error {
	result, err := DB.Exec(`UPDATE users SET is_active = ? WHERE id != 1`, isActive)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}

func UpdateActiveStatus(id string) error {
	var u User
	row := DB.QueryRow(`SELECT is_active FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No account found with the id %s", id)
		}
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	isActive := !u.IsActive
	_, err := DB.Exec(`UPDATE users SET is_active = ? WHERE id = ?`, isActive, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}

func GetAdminUsersCountAll() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return count, nil
}

func GetAdminUsersSearchCount(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam, from, to string) (int, error) {
	var count int
	var args []interface{}
	query := `SELECT COUNT(*) FROM users WHERE 
    (email LIKE ? OR name LIKE ? OR surname LIKE ? OR os LIKE ? OR created_at LIKE ?)`
	args = append(args, searchParam, searchParam, searchParam, searchParam, searchParam)

	if verifiedParam != "" {
		query += ` AND verified = ?`
		verified := verifiedParam == "1"
		args = append(args, verified)
	}

	if isActiveParam != "" {
		query += ` AND is_active = ?`
		isActive := isActiveParam == "1"
		args = append(args, isActive)
	}

	if isAdminParam != "" {
		query += ` AND is_admin = ?`
		isAdmin := isAdminParam == "1"
		args = append(args, isAdmin)
	}

	if specialAppsParam != "" {
		query += ` AND special_apps = ?`
		specialApps := specialAppsParam == "1"
		args = append(args, specialApps)
	}

	if from != "" && to != "" {
		query += ` AND created_at BETWEEN ? AND ?`
		args = append(args, from, to)
	}

	err := DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return count, nil
}

func GetAdminUsers(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam, from, to string, limit, cursor int) ([]User, error) {
	var users []User
	var args []interface{}
	query := `SELECT * FROM users WHERE 
    (email LIKE ? OR name LIKE ? OR surname LIKE ? OR os LIKE ? OR created_at LIKE ?)`

	args = append(args, searchParam, searchParam, searchParam, searchParam, searchParam)

	if verifiedParam != "" {
		query += ` AND verified = ?`
		verified := verifiedParam == "1"
		args = append(args, verified)
	}

	if isActiveParam != "" {
		query += ` AND is_active = ?`
		isActive := isActiveParam == "1"
		args = append(args, isActive)
	}

	if isAdminParam != "" {
		query += ` AND is_admin = ?`
		isAdmin := isAdminParam == "1"
		args = append(args, isAdmin)
	}

	if specialAppsParam != "" {
		query += ` AND special_apps = ?`
		specialApps := specialAppsParam == "1"
		args = append(args, specialApps)
	}

	if from != "" && to != "" {
		query += ` AND created_at BETWEEN ? AND ?`
		args = append(args, from, to)
	}

	query += ` ORDER BY ID DESC LIMIT ? OFFSET ?`
	args = append(args, limit, cursor)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
			&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
			&u.Pc, &u.Os, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return users, nil
}

func UpdatePassword(password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash the password: %v.", err)
	}

	result, err := DB.Exec(`UPDATE users SET password = ? WHERE email = ?`,
		hashedPassword, email)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v.", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v.", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found with the email %s.", email)
	}
	return nil
}

func DeleteAccountByEmail(email string) error {
	var userID int64
	err := DB.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No account found with the email %s", email)
		}
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	if userID == 1 {
		return fmt.Errorf("The account with ID 1 cannot be deleted")
	}

	result, err := DB.Exec(`DELETE FROM users WHERE email = ?`, email)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No account found with the email %s", email)
	}

	return nil
}

func VerifyAccount(userID int64) error {
	result, err := DB.Exec(`UPDATE users SET verified = true WHERE id = ?`, userID)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found with user id %d", userID)
	}
	return nil
}

func GetUserByEmail(email string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT * FROM users WHERE email = ?`, email)
	if err := row.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
		&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
		&u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("User not found with email %s", email)
		}
		return u, fmt.Errorf("Error finding user with email %s: %v", email, err)
	}
	return u, nil
}

func CreateUser(u User) error {
	date := utils.FormattedDate()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash the password: %v", err)
	}

	hashedPc, err := bcrypt.GenerateFromPassword([]byte(u.Pc), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash the unique identifier: %v", err)
	}

	_, err = DB.Exec(`
		INSERT INTO users
		(email, name, surname, password, is_admin, email_token, pc, os, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email, u.Name, u.Surname, string(hashedPassword), false, u.EmailToken, string(hashedPc), u.Os, date)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return fmt.Errorf("The email: %s already exists", u.Email)
		} else {
			return fmt.Errorf("An unexpected error occurred: %v", err)
		}
	}

	return nil
}
