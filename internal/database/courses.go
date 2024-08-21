package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type Course struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Thumbnail   string `json:"thumbnail"`
	Preview     string `json:"preview"`
	Rating      int    `json:"rating"`
	NumReviews  int    `json:"num_reviews"`
	Duration    string `json:"duration"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

func UpdateCourseActiveStatus(id string) error {
	// TODO: get the is_active from the frontend, !memory allocation
	var u User
	row := DB.QueryRow(`SELECT is_active FROM courses WHERE id = ?`, id)
	if err := row.Scan(&u.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: no such course", id)
		}
		return fmt.Errorf("%s: %v", id, err)
	}

	isActive := !u.IsActive
	_, err := DB.Exec(`UPDATE courses SET is_active = ? WHERE id = ?`, isActive, id)
	if err != nil {
		return fmt.Errorf("UpdateCourseActiveStatus: %v", err)
	}

	return nil
}

func GetCourseById(id string) (Course, error) {
	var c Course
	row := DB.QueryRow(`SELECT * FROM courses WHERE id = ?`, id)
	if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Author,
		&c.Thumbnail, &c.Preview, &c.Rating, &c.NumReviews, &c.Duration,
		&c.IsActive, &c.SortOrder, &c.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("GetCourseById: %s: no such course", id)
		}
		return c, fmt.Errorf("GetCourseById: %s: %v", id, err)
	}
	return c, nil
}

func DeleteCourseByID(id string) error {
	result, err := DB.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteCourseByID: course id: %s: %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteCourseByID: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("DeleteCourseByID: no course found with ID: %s", id)
	}
	return nil
}

func UpdateCourse(c Course) error {
	result, err := DB.Exec(`UPDATE courses SET 
  title = ?, description = ? , author = ?, thumbnail = ?, preview = ?, 
  duration = ?, is_active = ?, sort_order = ? WHERE id = ?`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration,
		c.IsActive, c.SortOrder, c.ID)
	if err != nil {
		return fmt.Errorf("UpdateCourse: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateCourse: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateCourse: no course found with ID: %d", c.ID)
	}
	return nil
}

func GetCourses(isActive string, searchTerm string) ([]Course, error) {
	var courses []Course
	query := `
		SELECT * FROM courses 
		WHERE (title LIKE ? OR description LIKE ? OR author LIKE ? OR duration LIKE ?)
	`
	args := []interface{}{
		searchTerm, searchTerm, searchTerm, searchTerm,
	}

	if isActive != "" {
		query += ` AND is_active = ?`
		args = append(args, isActive)
	}

	query += ` ORDER BY sort_order DESC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.Thumbnail,
			&c.Preview, &c.Rating, &c.NumReviews, &c.Duration, &c.IsActive,
			&c.SortOrder, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", err)
	}

	return courses, nil
}

func CreateCourse(c Course) (int64, error) {
	date := utils.FormattedDate()
	tx, err := DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var maxSortOrder int
	err = tx.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM courses").Scan(&maxSortOrder)
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: failed to get max sort_order: %v", err)
	}

	result, err := tx.Exec(`
        INSERT INTO courses
        (title, description, author, thumbnail, preview, duration, is_active, sort_order, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration, c.IsActive, maxSortOrder+1, date)
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: failed to insert course: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: failed to get last insert ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("CreateCourse: failed to commit transaction: %v", err)
	}

	return id, nil
}

func EditSortCourses(id int64, sort string) error {
	_, err := DB.Exec("UPDATE courses SET sort_order = ? WHERE id = ?", sort, id)
	if err != nil {
		return fmt.Errorf("EditSortCourses: %v", err)
	}
	return nil
}
