package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type Video struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoHLS    string `json:"video_hls"`
	Thumbnail   string `json:"thumbnail"`
	Length      string `json:"length"`
	Duration    string `json:"duration"`
	Views       int    `json:"views"`
	CourseID    int64  `json:"course_id"`
	CreatedAt   string `json:"created_at"`
	// not in db
	VideoResume string `json:"video_resume"`
}

// GetFeed obtiene videos relacionados a un curso con el último resumen, y permite búsqueda por título y descripción.
func GetFeed(userId int64, courseId, searchParam string) ([]Video, error) {
	query := `
		SELECT
			v.id AS video_id,
			v.title,
			v.description,
			v.video_hls,
			v.thumbnail,
			v.duration,
			v.length,
			v.views,
			v.course_id,
			COALESCE(h.video_resume, '') AS video_resume
		FROM
			videos v
		LEFT JOIN
			history h
			ON v.id = h.video_id
			AND h.user_id = ?
			AND h.created_at = (
				SELECT MAX(created_at)
				FROM history
				WHERE video_id = v.id
				AND user_id = ?
			)
		WHERE
			v.course_id = ? 
			AND (v.title LIKE ? OR v.description LIKE ?)
		ORDER BY
			v.id;
	`

	rows, err := DB.Query(query, userId, userId, courseId, "%"+searchParam+"%", "%"+searchParam+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var video Video
		err := rows.Scan(
			&video.ID,
			&video.Title,
			&video.Description,
			&video.VideoHLS,
			&video.Thumbnail,
			&video.Duration,
			&video.Length,
			&video.Views,
			&video.CourseID,
			&video.VideoResume,
		)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

func GetVideoById(videoId int64) (Video, error) {
	var v Video
	row := DB.QueryRow(`SELECT * FROM videos WHERE id = ?`, videoId)
	if err := row.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS,
		&v.Thumbnail, &v.Length, &v.Duration, &v.Views, &v.CourseID, &v.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return v, fmt.Errorf("GetVideoById: %d: no such video", videoId)
		}
		return v, fmt.Errorf("GetVideoById: %d: %v", videoId, err)
	}
	return v, nil
}

func GetFistVideoByCourseId(courseID string) (Video, error) {
	var video Video
	row := DB.QueryRow(`SELECT * FROM videos WHERE 
  course_id = ? ORDER BY id ASC LIMIT 1`, courseID)
	err := row.Scan(&video.ID, &video.Title, &video.Description,
		&video.VideoHLS, &video.Thumbnail, &video.Length, &video.Duration,
		&video.Views, &video.CourseID, &video.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return video, fmt.Errorf("GetFistVideoByCourseId: no video found with course ID: %s", courseID)
		}
		return video, fmt.Errorf("GetFistVideoByCourseId: error fetching first video: %v", err)
	}
	return video, nil
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

func GetVideosCount(course_id, searchParam string) (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM videos 
  WHERE course_id = ? AND (title LIKE ? OR description LIKE ?)`,
		course_id, "%"+searchParam+"%", "%"+searchParam+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetVideosCount: %v", err)
	}
	return count, nil
}

func GetAdminVideos(course_id string, searchParam string, offset int, limit int) ([]Video, error) {
	var videos []Video
	rows, err := DB.Query(`SELECT * FROM videos
		WHERE course_id = ? AND (title LIKE ? OR description LIKE ?)
		ORDER BY id DESC
		LIMIT ? OFFSET ?`,
		course_id, "%"+searchParam+"%", "%"+searchParam+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetAdminVideos: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS,
			&v.Thumbnail, &v.Duration, &v.Length, &v.Views, &v.CourseID,
			&v.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAdminVideos: %v", err)
		}
		videos = append(videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAdminVideos: %v", err)
	}
	return videos, nil
}

func GetVideos(course_id string, searchParam string, offset int, limit int) ([]Video, error) {
	var videos []Video
	rows, err := DB.Query(`SELECT * FROM videos
		WHERE course_id = ? AND (title LIKE ? OR description LIKE ?)
		ORDER BY id
		LIMIT ? OFFSET ?`,
		course_id, "%"+searchParam+"%", "%"+searchParam+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetVideos: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS,
			&v.Thumbnail, &v.Duration, &v.Length, &v.Views, &v.CourseID,
			&v.CreatedAt); err != nil {
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
	date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO videos
  (title, description, video_hls, thumbnail, length, duration, course_id, created_at) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		v.Title, v.Description, v.VideoHLS, v.Thumbnail, v.Length, v.Duration, v.CourseID, date)

	if err != nil {
		return 0, fmt.Errorf("CreateVideo: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateVideo: %v", err)
	}

	return id, nil
}
