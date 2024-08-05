package database

import "fmt"

type History struct {
  ID            int64   `json:"id"`
  VideoId       int64   `json:"video_id"`
  CourseId      int64   `json:"course_id"`
  UserId        int64   `json:"user_id"`
  VideoResume        string  `json:"video_resume"`
  CreatedAt     string  `json:"created_at"`
}

func CreateHistory(user_id int64, video_id string, course_id string, resume string) (History, error) {
  result, err := DB.Exec(`
  INSERT INTO history
  (user_id, video_id, course_id, video_resume) 
  VALUES (?, ?, ?, ?)`, 
  user_id, video_id, course_id, resume)

	if err != nil {
    return History{}, fmt.Errorf("CreateHistory: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
    return History{}, fmt.Errorf("CreateHistory: %v", err)
	}

	insertedHistory := History{
		ID:         id,
    UserId:     user_id,
    VideoResume:     resume,
	}
	return insertedHistory, nil
}

func UpdateHistory(id string, resume string) error {
  _, err := DB.Exec(`UPDATE history SET video_resume = ? WHERE id = ?`, resume, id)
  if err != nil {
    return fmt.Errorf("UpdateHistory: %v", err)
  }
  return nil
}

func GetLastVideoByUserIdAndCourseIdAndVideoId(user_id int64, course_id string, video_id string) (History, error) {
  var history History
  row := DB.QueryRow(`
  SELECT 
  *
  FROM history 
  WHERE user_id = ? AND course_id = ? AND video_id = ?
  ORDER BY created_at DESC
  LIMIT 1`, user_id, course_id, video_id)

  err := row.Scan(&history.ID, &history.VideoId, &history.CourseId, &history.UserId, &history.VideoResume, &history.CreatedAt)
  if err != nil {
    return history, fmt.Errorf("GetHistoryByUserIdAndCourseId: %v", err)
  }

  return history, nil
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
    return history, fmt.Errorf("GetHistoryByUserIdAndCourseId: %v", err)
  }

  return history, nil
}

func GetUserHistory(user_id int64) ([]History, error) {
  var histories []History
  rows, err := DB.Query(`SELECT 
    *
    FROM history WHERE user_id = ?`, user_id)
  if err != nil {
    return histories, fmt.Errorf("GetUserHistory: %v", err)
  }
  defer rows.Close()
  for rows.Next() {
    var history History
    if err := rows.Scan(&history.ID, &history.VideoId, &history.CourseId, &history.UserId, &history.VideoResume, &history.CreatedAt); err != nil {
      return histories, fmt.Errorf("GetUserHistory: %v", err)
    }
    histories = append(histories, history)
  }
  return histories, nil
}
