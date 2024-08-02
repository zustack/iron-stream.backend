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

func CreateHistory(h History) (int64, error) {
  result, err := DB.Exec(`
  INSERT INTO history
  (user_id, video_id, course_id, video_resume) 
  VALUES (?, ?, ?, ?)`, 
  h.UserId, h.VideoId, h.CourseId, h.VideoResume)

	if err != nil {
    return 0, fmt.Errorf("CreateHistory: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
    return 0, fmt.Errorf("CreateHistory: %v", err)
	}

	return id, nil
}

func GetHistoryByUserID(userID int64) ([]History, error) {
  var history []History
  rows, err := DB.Query(`SELECT * FROM history WHERE user_id = ?`, userID)
  if err != nil {
    return nil, fmt.Errorf("GetHistoryByUserID: %v", err)
  }
  defer rows.Close()

  for rows.Next() {
    var h History
    if err := rows.Scan(&h.ID, &h.VideoId, &h.CourseId, &h.UserId, &h.VideoResume, &h.CreatedAt); err != nil {
      return nil, fmt.Errorf("GetHistoryByUserID: %v", err)
    }
    history = append(history, h)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("GetHistoryByUserID: %v", err)
  }

  return history, nil
}

func GetLastVideoRecord(user_id int64) (History, error) {
  var h History
  row := DB.QueryRow(`SELECT * FROM history WHERE 
  user_id = ? ORDER BY created_at DESC LIMIT 1`, user_id)
  err := row.Scan(&h.ID, &h.VideoId, &h.CourseId, &h.UserId, &h.VideoResume, &h.CreatedAt)
  if err != nil {
    return h, fmt.Errorf("GetLastVideoRecord: %v", err)
  }
  return h, nil
}

func GetHistoryCountByUserID(userID int64) (int, error) {
  var count int
  err := DB.QueryRow("SELECT COUNT(*) FROM history WHERE user_id = ?", userID).Scan(&count)
  if err != nil {
    return 0, fmt.Errorf("GetHistoryCountByUserID: %v", err)
  }
  return count, nil
}

func UpdateHistory(id int64, video_resume string) error {
  result, err := DB.Exec(`UPDATE history SET video_resume = ? WHERE id = ?`, video_resume, id)
  if err != nil {
    return fmt.Errorf("UpdateHistory: %v", err)
  }
  _, err = result.RowsAffected()
  if err != nil {
    return fmt.Errorf("UpdateHistory: %v", err)
  }
  return nil
}
