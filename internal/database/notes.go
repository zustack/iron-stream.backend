package database

import "fmt"

type Note struct {
	ID             int64  `json:"id"`
  Body           string `json:"body"`
  VideoTitle     string `json:"video_title"`
  Time           string `json:"time"`
  CourseID       int64  `json:"course_id"`
  UserID         int64  `json:"user_id"`
}

func DeleteNoteById(id string) error {
	result, err := DB.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No note found with the id %s", id)
	}
	return nil
}

func UpdateNote(id, body string) error {
	result, err := DB.Exec(`UPDATE notes SET body = ? WHERE id = ?`,
		body, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No note found with the id %v", id)
	}
	return nil
}

func GetNotesByCourseIdAndUserId(courseId string, userId int64) ([]Note, error) {
	var notes []Note
	rows, err := DB.Query(`SELECT * FROM notes WHERE course_id = ? AND user_id = ?`, courseId, userId)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Body, &n.VideoTitle, &n.Time, &n.CourseID, &n.UserID); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return notes, nil
}

func CreateNote(n Note) error {
	_, err := DB.Exec(`
  INSERT INTO apps
  (body, video_title, time, course_id, user_id) 
  VALUES (?, ?, ?, ?, ?)`,
    n.Body, n.VideoTitle, n.Time, n.CourseID, n.UserID)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
