package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

func AddCourseToUser(userID, courseID int64) error {
	var courseExists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)", courseID).Scan(&courseExists)
	if err != nil {
		return fmt.Errorf("CheckCourseExists: %v", err)
	}
	if !courseExists {
		return fmt.Errorf("Course with ID %d does not exist", courseID)
	}

	var existingCourses string
	row := DB.QueryRow(`SELECT courses FROM users WHERE id = ?`, userID)
	if err := row.Scan(&existingCourses); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("GetUserByID %d: no such user", userID)
		}
		return fmt.Errorf("GetUserByID %d: %v", userID, err)
	}

	courseList := strings.Split(existingCourses, ",")
	for _, b := range courseList {
		if b == strconv.FormatInt(courseID, 10) {
			return nil
		}
	}

	var coursesIDs string
	if existingCourses == "" {
		coursesIDs = strconv.FormatInt(courseID, 10)
	} else {
		coursesIDs = existingCourses + "," + strconv.FormatInt(courseID, 10)
	}

	_, err = DB.Exec("UPDATE users SET courses = ? WHERE id = ?", coursesIDs, userID)
	if err != nil {
		return fmt.Errorf("AddCourseToUser: %v", err)
	}

	return nil
}

func GetCourseClientCount(searchParam, isActiveParam string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM courses
		WHERE (title LIKE ? OR description LIKE ? OR author LIKE ? OR duration LIKE ?)
		AND is_active = ?`

	searchPattern := "%" + searchParam + "%"

	err := DB.QueryRow(query, searchPattern, searchPattern, searchPattern, searchPattern, isActiveParam).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetCourseClientCount: %v", err)
	}
	return count, nil
}

func GetAdminCourses(searchParam, isActiveParam string, limit, cursor int) ([]Course, error) {
	var courses []Course
	var args []interface{}
	query := `SELECT * FROM courses WHERE 
              (title LIKE ? OR description LIKE ? OR author LIKE ? OR duration LIKE ?)`

	args = append(args, searchParam, searchParam, searchParam, searchParam)

	if isActiveParam != "" {
		query += ` AND is_active = ?`
		isActive := isActiveParam == "1"
		args = append(args, isActive)
	}

	query += ` ORDER BY sort_order DESC LIMIT ? OFFSET ?`
	args = append(args, limit, cursor)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAdminCourses: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author,
			&c.Thumbnail, &c.Preview, &c.Rating, &c.NumReviews,
			&c.Duration, &c.IsActive, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAdminCourses: %v", err)
		}
		courses = append(courses, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAdminCourses: %v", err)
	}
	return courses, nil
}

func GetCoursesCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM courses").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetCoursesCount: %v", err)
	}
	return count, nil
}

func GetCourses(searchParam string, offset int, limit int) ([]Course, error) {
	var courses []Course
	rows, err := DB.Query(`SELECT * FROM courses
  WHERE is_active = 1 AND title LIKE ? OR description LIKE ? OR author LIKE ? OR duration LIKE ? 
  ORDER BY sort_order DESC LIMIT ? OFFSET ?`,
		searchParam, searchParam, searchParam, searchParam, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetCourses: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.Thumbnail, &c.Preview, &c.Rating, &c.NumReviews, &c.Duration, &c.IsActive, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCourses: %v", err)
		}
		courses = append(courses, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCourses: %v", err)
	}
	return courses, nil
}

func CreateCourse(c Course) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO courses
  (title, description, author, thumbnail, preview, duration, is_active, sort_order) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration, c.IsActive, 0)

	if err != nil {
		return 0, fmt.Errorf("CreateCourse: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: %v", err)
	}

	_, err = DB.Exec("UPDATE courses SET sort_order = ? WHERE id = ?", id, id)
	if err != nil {
		return 0, fmt.Errorf("CreateCourse/updateOrder: %v", err)
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
