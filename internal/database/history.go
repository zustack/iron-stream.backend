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

type HistoryWithVideo struct {
	HistoryID    int64  `json:"history_id"`
	VideoID      int64  `json:"video_id"`
	CourseID     int64  `json:"course_id"`
	UserID       int64  `json:"user_id"`
	VideoResume  string `json:"video_resume"`
	HistoryDate  string `json:"history_date"`
	VideoTitle   string `json:"video_title"`
	Description  string `json:"description"`
	VideoHLS     string `json:"video_hls"`
	Thumbnail    string `json:"thumbnail"`
	Duration     string `json:"duration"`
	Length       string `json:"length"`
	Views        int64  `json:"views"`
	VideoCreated string `json:"video_created"`
}

func GetUserHistoryWithVideos(userID int64) ([]HistoryWithVideo, error) {
	var userHistory []HistoryWithVideo

	query := `
SELECT h.id, h.video_id, h.course_id, h.user_id, h.video_resume, h.created_at, 
       v.title, v.description, v.video_hls, v.thumbnail, v.duration, v.length, v.views, v.created_at
FROM history h
INNER JOIN videos v ON h.video_id = v.id
INNER JOIN (
    SELECT video_id, MAX(created_at) as latest_history
    FROM history
    WHERE user_id = ?
    GROUP BY video_id
) latest ON h.video_id = latest.video_id AND h.created_at = latest.latest_history
WHERE h.user_id = ?
ORDER BY h.created_at DESC;
`

	rows, err := DB.Query(query, userID, userID)
	if err != nil {
		return userHistory, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var hv HistoryWithVideo
		if err := rows.Scan(&hv.HistoryID, &hv.VideoID, &hv.CourseID, &hv.UserID, &hv.VideoResume, &hv.HistoryDate,
			&hv.VideoTitle, &hv.Description, &hv.VideoHLS, &hv.Thumbnail, &hv.Duration, &hv.Length, &hv.Views, &hv.VideoCreated); err != nil {
			return userHistory, fmt.Errorf("An unexpected error occurred while scanning rows: %v", err)
		}
		userHistory = append(userHistory, hv)
	}

	if err := rows.Err(); err != nil {
		return userHistory, fmt.Errorf("An unexpected error occurred with rows: %v", err)
	}

	return userHistory, nil
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
