package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type Course struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Author         string `json:"author"`
	Thumbnail      string `json:"thumbnail"`
	Preview        string `json:"preview"`
	Rating         int    `json:"rating"`
	NumReviews     int    `json:"num_reviews"`
	Duration       string `json:"duration"`
	IsActive       bool   `json:"is_active"`
  Price int `json:"price"`
	SortOrder      int    `json:"sort_order"`
	CreatedAt      string `json:"created_at"`
	IsUserEnrolled bool   `json:"is_user_enrolled"`
}

func GetCourses(isActive string, searchTerm string) ([]Course, error) {
	var courses []Course
	query := `
		SELECT * FROM courses 
		WHERE (title LIKE ? OR description LIKE ? OR author LIKE ? OR 
    duration LIKE ?)
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
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.Thumbnail,
			&c.Preview, &c.Rating, &c.NumReviews, &c.Duration, &c.IsActive, &c.Price,
			&c.SortOrder, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return courses, nil
}

func UpdateCourseActiveStatus(id string, active string) error {
	result, err := DB.Exec(`UPDATE courses SET is_active = ? WHERE id = ?`, active, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No course found with the id %v", id)
	}
	return nil
}

func EditSortCourses(id int64, sort string) error {
	result, err := DB.Exec("UPDATE courses SET sort_order = ? WHERE id = ?", sort, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No course found with the id %v", id)
	}
	return nil
}

func GetCourseById(id string) (Course, error) {
	var c Course
	row := DB.QueryRow(`SELECT * FROM courses WHERE id = ?`, id)
	if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Author,
		&c.Thumbnail, &c.Preview, &c.Rating, &c.NumReviews, &c.Duration,
		&c.IsActive, &c.Price, &c.SortOrder, &c.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("No course found with the id %s", id)
		}
		return c, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return c, nil
}

func DeleteCourseByID(id string) error {
	result, err := DB.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No course found with the id %v", id)
	}
	return nil
}

func UpdateCourse(c Course) error {
	result, err := DB.Exec(`UPDATE courses SET 
  title = ?, description = ? , author = ?, thumbnail = ?, preview = ?, 
  duration = ?, is_active = ? WHERE id = ?`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration,
		c.IsActive, c.ID)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No course found with the id %v", c.ID)
	}
	return nil
}

func CreateCourse(c Course) error {
	var maxSortOrder int
	err := DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM courses").Scan(&maxSortOrder)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	date := utils.FormattedDate()
	_, err = DB.Exec(`
        INSERT INTO courses
        (title, description, author, thumbnail, preview, duration, is_active, sort_order, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration, c.IsActive, maxSortOrder+1, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
