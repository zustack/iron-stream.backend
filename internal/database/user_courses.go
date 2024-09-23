package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type CourseProfit struct {
	Title  string `json:"course"`
	Profit int    `json:"profit"`
}

type CourseProfitResult struct {
	Courses []CourseProfit `json:"courses"`
	Total   int            `json:"total"`
}

func GetCoursesProfit(from, to string) (CourseProfitResult, error) {
	query := `
SELECT c.title, SUM(uc.price) as profit
FROM courses c
JOIN user_courses uc ON c.id = uc.course_id
WHERE uc.created_at BETWEEN ? AND ? AND uc.user_id != 1
GROUP BY c.id, c.title
ORDER BY profit DESC
	`
	rows, err := DB.Query(query, from, to)
	if err != nil {
		return CourseProfitResult{}, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var result CourseProfitResult
	var total int

	for rows.Next() {
		var cp CourseProfit
		if err := rows.Scan(&cp.Title, &cp.Profit); err != nil {
			return CourseProfitResult{}, fmt.Errorf("failed to scan row: %v", err)
		}
		result.Courses = append(result.Courses, cp)
		total += cp.Profit
	}

	if err := rows.Err(); err != nil {
		return CourseProfitResult{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	result.Total = total
	return result, nil
}

func DeleteUserCoursesByCourseIdAndUserId(userId string, courseId string) error {
	result, err := DB.Exec(`DELETE FROM user_courses WHERE course_id = ? AND user_id = ?;`, courseId, userId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func DeleteUserCoursesByCourseId(courseId string) error {
	result, err := DB.Exec(`DELETE FROM user_courses WHERE course_id = ?;`, courseId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No course found with the id %v", courseId)
	}
	return nil
}

func DeleteAllUserCourses() error {
	result, err := DB.Exec(`DELETE FROM user_courses;`)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func UserCourseExists(userID int64, courseID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM user_courses 
			WHERE user_id = ? AND course_id = ?
		);`
	var exists bool
	err := DB.QueryRow(query, userID, courseID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return exists, nil
}

func GetUserCourses(userID, q string) ([]Course, error) {
	query := `
		SELECT 
			c.id,
			c.title,
			c.description,
			c.author,
			c.thumbnail,
			c.preview,
			c.rating,
			c.num_reviews,
			c.duration,
			c.is_active,
			c.sort_order,
			c.created_at,
			CASE 
				WHEN uc.user_id IS NOT NULL THEN true 
				ELSE false 
			END AS is_user_enrolled
		FROM 
			courses c
		LEFT JOIN 
			user_courses uc 
		ON 
			c.id = uc.course_id 
			AND uc.user_id = ?
		WHERE 
		  LOWER(c.title) LIKE '%' || LOWER(?) || '%' 
			OR LOWER(c.author) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.duration) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.created_at) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.description) LIKE '%' || LOWER(?) || '%';`

	rows, err := DB.Query(query, userID, q, q, q, q, q)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Author,
			&course.Thumbnail,
			&course.Preview,
			&course.Rating,
			&course.NumReviews,
			&course.Duration,
			&course.IsActive,
			&course.SortOrder,
			&course.CreatedAt,
			&course.IsUserEnrolled,
		)
		if err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return courses, nil
}

func CreateUserCourse(userID, courseID string) error {
	var userExists bool
	err := DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`,
		userID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if !userExists {
		return fmt.Errorf("The user with the id %s does not exist", userID)
	}

	var price int
	row := DB.QueryRow(`SELECT price FROM courses WHERE id = ?`, courseID)
	if err := row.Scan(&price); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("The course with the id %s does not exist", courseID)
		}
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	date := utils.FormattedDate()

	_, err = DB.Exec(`
    INSERT INTO user_courses (user_id, course_id, created_at, price) VALUES (?, ?, ?, ?)`,
		userID, courseID, date, price)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
