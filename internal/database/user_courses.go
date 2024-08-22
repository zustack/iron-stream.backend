package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

func DeleteUserCoursesByCourseIdAndUserId(userId string, courseId string) error {
  fmt.Println("courseId", courseId)
  fmt.Println("userId", userId)
	_, err := DB.Exec(`DELETE FROM user_courses WHERE course_id = ? AND user_id = ?;`, courseId, userId)
	if err != nil {
		return fmt.Errorf("DeleteUserCoursesByUserID: %v", err)
	}
	return nil
}


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

func CreateUserCourse(userID, courseID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("CreateUserCourse: %v", err)
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
		return fmt.Errorf("CreateUserCourse: %v", err)
	}
	if !userExists {
		return fmt.Errorf("CreateUserCourse: userID %s no existe", userID)
	}

	var courseExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)`,
		courseID).Scan(&courseExists)
	if err != nil {
		return fmt.Errorf("CreateUserCourse: %v", err)
	}
	if !courseExists {
		return fmt.Errorf("CreateUserCourse: courseID %s no existe", courseID)
	}

	date := utils.FormattedDate()
	_, err = tx.Exec(`
    INSERT INTO user_courses (user_id, course_id, created_at) VALUES (?, ?, ?)`,
		userID, courseID, date)
	if err != nil {
		return fmt.Errorf("CreateUserCourse: %v", err)
	}

	return nil
}
