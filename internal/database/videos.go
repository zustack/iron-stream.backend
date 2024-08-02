package database

import (
	"database/sql"
	"fmt"
)

type Video struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoHLS    string `json:"video_hls"`
	Thumbnail   string `json:"thumbnail"`
	Length      string `json:"length"`
	Views       int    `json:"views"`
	CourseID    int64  `json:"course_id"`
	SortOrder   int64  `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

func UpdateVideoViews(id int64) error {
  // First, get the current view count
  var currentViews int
  err := DB.QueryRow("SELECT views FROM videos WHERE id = ?", id).Scan(&currentViews)
  if err != nil {
    if err == sql.ErrNoRows {
      return fmt.Errorf("UpdateVideoViews: no video found with ID: %d", id)
    }
    return fmt.Errorf("UpdateVideoViews: error fetching current views: %v", err)
  }

  // Increment the view count
  newViews := currentViews + 1

  // Update the database with the new view count
  result, err := DB.Exec(`UPDATE videos SET views = ? WHERE id = ?`, newViews, id)
  if err != nil {
    return fmt.Errorf("UpdateVideoViews: %v", err)
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("UpdateVideoViews: error getting rows affected: %v", err)
  }

  if rowsAffected == 0 {
    return fmt.Errorf("UpdateVideoViews: no video updated with ID: %d", id)
  }

  return nil
}

func DeleteVideoByID(id string) error {
	result, err := DB.Exec("DELETE FROM videos WHERE id = ?", id)
	if err != nil {
    return fmt.Errorf("DeleteVideoByID: course id: %s: %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
    return fmt.Errorf("DeleteVideoByID: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
    return fmt.Errorf("DeleteVideoByID: no video found with ID: %s", id)
	}
	return nil
}

func UpdateVideo(v Video) error {
	result, err := DB.Exec(`UPDATE videos SET 
  title = ?, description = ?, video_hls = ?, thumbnail = ?, length = ?
  WHERE id = ?`,
  v.Title, v.Description, v.VideoHLS, v.Thumbnail, v.Length, v.ID)
	if err != nil {
		return fmt.Errorf("UpdateVideo: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateVideo: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateVideo: no video found with ID: %d", v.ID)
	}
	return nil
}

func GetVideosCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetVideosCount: %v", err)
	}
	return count, nil
}

func GetVideos(course_id, searchParam string, offset int, limit int) ([]Video, error) {
	var videos []Video
	rows, err := DB.Query(`SELECT * FROM videos
  WHERE course_id = ? AND title LIKE ? OR description LIKE ? 
  ORDER BY sort_order LIMIT ? OFFSET ?`,
		course_id, searchParam, searchParam, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetVideos: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS, 
    &v.Thumbnail, &v.Length, &v.Views, &v.CourseID, &v.SortOrder, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetVideos: %v", err)
		}
		videos = append(videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetVideos: %v", err)
	}
	return videos, nil
}

func CreateVideo(v Video) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO videos
  (title, description, video_hls, thumbnail, length, course_id, sort_order) 
  VALUES (?, ?, ?, ?, ?, ?, ?)`,
		v.Title, v.Description, v.VideoHLS, v.Thumbnail, v.Length, v.CourseID, 0)

	if err != nil {
		return 0, fmt.Errorf("CreateVideo: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateVideo: %v", err)
	}

	_, err = DB.Exec("UPDATE videos SET sort_order = ? WHERE id = ?", id, id)
	if err != nil {
		return 0, fmt.Errorf("CreateVideo/updateSortOrder: %v", err)
	}

	return id, nil
}
