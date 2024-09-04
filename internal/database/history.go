package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type History struct {
	ID          int64  `json:"id"`
	VideoId     string `json:"video_id"`
	CourseId    int64  `json:"course_id"`
	UserId      int64  `json:"user_id"`
	VideoResume string `json:"video_resume"`
	CreatedAt   string `json:"created_at"`
}

func GetLastVideoByUserIdAndCourseIdAndVideoId(userId int64, courseId string, videoId int64) (History, error) {
	var history History
	row := DB.QueryRow(`SELECT * FROM history 
  WHERE user_id = ? AND course_id = ? AND video_id = ?
  ORDER BY created_at DESC LIMIT 1`, userId, courseId, videoId)
	err := row.Scan(&history.ID, &history.VideoId, &history.CourseId,
		&history.UserId, &history.VideoResume, &history.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return history, fmt.Errorf("Record not found")
		}
		return history, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return history, nil
}

func UpdateHistory(id string, resume string) error {
	result, err := DB.Exec(`UPDATE history SET video_resume = ? WHERE id = ?`, resume, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No history with id: %s", id)
	}
	return nil
}

func GetLastVideoByUserIdAndCourseId(user_id int64, course_id string) (History, error) {
	var history History
	row := DB.QueryRow(`
  SELECT 
  *
  FROM history 
  WHERE user_id = ? AND course_id = ?
  ORDER BY created_at DESC
  LIMIT 1`, user_id, course_id)

	err := row.Scan(&history.ID, &history.VideoId, &history.CourseId, &history.UserId, &history.VideoResume, &history.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return history, fmt.Errorf("Record not found")
		}
		return history, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return history, nil
}

func GetUserHistory(user_id int64) ([]History, error) {
	var hs []History
	rows, err := DB.Query(`SELECT * FROM history WHERE user_id = ?`, user_id)
	if err != nil {
		return hs, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.VideoId, &h.CourseId, &h.UserId, &h.VideoResume, &h.CreatedAt); err != nil {
			return hs, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func CreateHistory(userId, videoId int64, courseId string, resume string) (History, error) {
	date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO history
  (user_id, video_id, course_id, video_resume, created_at) 
  VALUES (?, ?, ?, ?, ?)`,
		userId, videoId, courseId, resume, date)

	if err != nil {
		return History{}, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return History{}, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	insertedHistory := History{
		ID:          id,
		VideoResume: resume,
	}

	return insertedHistory, nil
}
