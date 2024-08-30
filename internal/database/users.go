package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
	"strconv"
	"strings"
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

func UpdateAdminStatus(userId, isAdmin string) error {
	_, err := DB.Exec(`UPDATE users SET is_admin = ? WHERE id = ?`, isAdmin, userId)
	if err != nil {
		return fmt.Errorf("UpdateAdminStatus: %v", err)
	}
	return nil
}

func UpdateUserSpecialApps(userId, special_apps string) error {
	_, err := DB.Exec(`UPDATE users SET special_apps = ? WHERE id = ?`, special_apps, userId)
	if err != nil {
		return fmt.Errorf("UpdateUserSpecialApps: %v", err)
	}
	return nil
}

// change
func DeactivateCourseForUser(userID, courseID int64) error {
	var courseExists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)", courseID).Scan(&courseExists)
	if err != nil {
		return fmt.Errorf("CheckCourseExists: %v", err)
	}
	if !courseExists {
		return fmt.Errorf("Course with ID %d does not exist", courseID)
	}

	// Obtiene la lista de cursos del usuario
	var existingCourses string
	row := DB.QueryRow(`SELECT courses FROM users WHERE id = ?`, userID)
	if err := row.Scan(&existingCourses); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("GetUserByID %d: no such user", userID)
		}
		return fmt.Errorf("GetUserByID %d: %v", userID, err)
	}

	// Divide la lista de cursos en un slice
	courseList := strings.Split(existingCourses, ",")
	newCourseList := make([]string, 0, len(courseList))

	// Elimina el courseID de la lista de cursos
	// courseFound := false
	for _, b := range courseList {
		if b != strconv.FormatInt(courseID, 10) {
			newCourseList = append(newCourseList, b)
		} else {
			// courseFound = true
		}
	}

	/*
		if !courseFound {
			return fmt.Errorf("Course with ID %d is not enrolled by user %d", courseID, userID)
		}
	*/

	// Une la lista de cursos actualizada
	updatedCourses := strings.Join(newCourseList, ",")

	// Actualiza los cursos del usuario en la base de datos
	_, err = DB.Exec("UPDATE users SET courses = ? WHERE id = ?", updatedCourses, userID)
	if err != nil {
		return fmt.Errorf("DeactivateCourseForUser: %v", err)
	}

	return nil
}

func GetUserIds() ([]User, error) {
	var ids []User
	rows, err := DB.Query(`SELECT ID FROM users`)
	if err != nil {
		return nil, fmt.Errorf("GetUserIds: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID); err != nil {
			return nil, fmt.Errorf("GetUserIds: %v", err)
		}
		ids = append(ids, i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserIds: %v", err)
	}
	return ids, nil
}

// change
func DeactivateAllCoursesForAllUsers() error {
	_, err := DB.Exec("UPDATE users SET courses = ''")
	if err != nil {
		return fmt.Errorf("DeactivateAllCoursesForAllUsers: %v", err)
	}
	return nil
}

func UpdateActiveStatusAllUsers(isActive bool) error {
	result, err := DB.Exec(`UPDATE users SET is_active = ? WHERE is_admin = false`, isActive)
	if err != nil {
		return fmt.Errorf("UpdateActiveStatusAllUsers: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateActiveStatusAllUsers: error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("UpdateActiveStatusAllUsers: no rows were updated")
	}

	return nil
}

func UpdateActiveStatus(id string) error {
	var u User
	row := DB.QueryRow(`SELECT is_active FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("GetUserByID %s: no such user", id)
		}
		return fmt.Errorf("GetUserByID %s: %v", id, err)
	}

	isActive := !u.IsActive
	_, err := DB.Exec(`UPDATE users SET is_active = ? WHERE id = ?`, isActive, id)
	if err != nil {
		return fmt.Errorf("UpdateActiveStatus: %v", err)
	}

	return nil
}

func GetAdminUsersCountAll() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetAdminUsersCountAll: %v", err)
	}
	return count, nil
}

func GetAdminUsersCount(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam string) (int, error) {
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
	err := DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetAdminUsersCount: %v", err)
	}
	return count, nil
}

func GetAdminUsers(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam string, limit, cursor int) ([]User, error) {
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

	query += ` ORDER BY ID DESC LIMIT ? OFFSET ?`
	args = append(args, limit, cursor)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAdminUsers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Password, &u.Email, &u.Name,
			&u.Surname, &u.IsAdmin, &u.SpecialApps, &u.IsActive, &u.EmailToken, &u.Verified,
			&u.Pc, &u.Os, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAdminUsers: %v", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAdminUsers: %v", err)
	}
	return users, nil
}

func UpdateEmailToken(email string, email_token int) error {
	_, err := DB.Exec(`UPDATE users SET email_token = ? WHERE email = ?`, email_token, email)
	if err != nil {
		return fmt.Errorf("UpdateEmailToken: %v", err)
	}
	return nil
}

func UpdatePassword(password, email string) error {
	email = strings.TrimSpace(email)
	result, err := DB.Exec(`UPDATE users SET password = ? WHERE LOWER(email) = LOWER(?)`, password, email)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %v", err)
	}
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
		&u.Pc, &u.Os, &u.CreatedAt); err != nil {
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
		&u.Pc, &u.Os, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("User not found with email %s", email)
		}
		return u, fmt.Errorf("Error finding user with email %s: %v", email, err)
	}
	return u, nil
}

func CreateUser(u User) (int64, error) {
	date := utils.FormattedDate()
	emailToken := utils.GenerateCode()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	result, err := DB.Exec(`
		INSERT INTO users
		(email, name, surname, password, is_admin, email_token, pc, os, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email, u.Name, u.Surname, hashedPassword, false, emailToken, u.Pc, u.Os, date)

	if err != nil {
		return 0, fmt.Errorf("The email: %s already exists", u.Email)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	return id, nil
}
