package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

func DeleteUserCoursesByCourseId(courseId string) error {
	_, err := DB.Exec(`DELETE FROM user_courses WHERE course_id = ?;`, courseId)
	if err != nil {
		return fmt.Errorf("DeleteUserCoursesByUserID: %v", err)
	}
	return nil
}

func DeleteAllUserCourses() error {
	_, err := DB.Exec(`DELETE FROM user_courses;`)
	if err != nil {
		return fmt.Errorf("DeleteAllUserCourses: %v", err)
	}
	return nil
}

func GetUserCourses(userID string) ([]Course, error) {
	var courses []Course
	rows, err := DB.Query(`SELECT c.* FROM courses c 
    JOIN user_courses uc ON c.id = uc.course_id 
    WHERE uc.user_id = ?;
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserCourses: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.Thumbnail,
			&c.Preview, &c.Rating, &c.NumReviews, &c.Duration, &c.IsActive,
			&c.SortOrder, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetUserCourses: %v", err)
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserCourses: %v", err)
	}

	return courses, nil
}

func AddCourseToUser(userID, courseID int64) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO user_courses (user_id, course_id, created_at) VALUES (?, ?, ?)`,
		userID, courseID, date)
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}
	return nil
}
