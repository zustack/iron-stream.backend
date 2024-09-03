package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

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

func GetUserCourseIds(userID string) ([]int64, error) {
	var courseIDs []int64
	rows, err := DB.Query(`
		SELECT uc.course_id FROM user_courses uc
		WHERE uc.user_id = ?;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var courseID int64
		if err := rows.Scan(&courseID); err != nil {
			return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
		}
		courseIDs = append(courseIDs, courseID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
	}

	return courseIDs, nil
}

// this is for the middleware of the videos
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
		return false, err
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
			AND uc.user_id = $1
		WHERE 
		  LOWER(c.title) LIKE '%' || LOWER(?) || '%' 
			OR LOWER(c.author) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.duration) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.created_at) LIKE '%' || LOWER(?) || '%'
			OR LOWER(c.description) LIKE '%' || LOWER(?) || '%';`

	rows, err := DB.Query(query, userID, q, q, q, q, q)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func CreateUserCourse(userID, courseID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`,
		userID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if !userExists {
		return fmt.Errorf("The user with the id %s does not exist", userID)
	}

	var courseExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)`,
		courseID).Scan(&courseExists)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if !courseExists {
		return fmt.Errorf("The course with the id %s does not exist", courseID)
	}

	date := utils.FormattedDate()
	_, err = tx.Exec(`
    INSERT INTO user_courses (user_id, course_id, created_at) VALUES (?, ?, ?)`,
		userID, courseID, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
