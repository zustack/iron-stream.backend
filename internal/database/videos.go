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
	CourseID    string `json:"course_id"`
	CreatedAt   string `json:"created_at"`
	VideoResume string `json:"video_resume"`
	SReview     bool   `json:"s_review"`
}

func UpdateVideoSReview(s_review, id string) error {
	result, err := DB.Exec(`UPDATE videos SET 
  s_review = ? WHERE id = ?`, s_review, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No video found with id %d", id)
	}
	return nil
}

func GetVideoById(videoId string) (Video, error) {
	var v Video
	row := DB.QueryRow(`SELECT * FROM videos WHERE id = ?`, videoId)
	if err := row.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS,
		&v.Thumbnail, &v.Length, &v.Duration, &v.Views, &v.CourseID, &v.CreatedAt, &v.SReview); err != nil {
		if err == sql.ErrNoRows {
			return v, fmt.Errorf("No video found with id %s", videoId)
		}
		return v, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return v, nil
}

func UpdateVideoViews(id int64) error {
	var currentViews int
	err := DB.QueryRow("SELECT views FROM videos WHERE id = ?", id).Scan(&currentViews)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No video found with id %v", id)
		}
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	newViews := currentViews + 1

	result, err := DB.Exec(`UPDATE videos SET views = ? WHERE id = ?`, newViews, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No video updated with ID: %d", id)
	}

	return nil
}

func GetFistVideoByCourseId(courseID string) (Video, error) {
	var video Video
	row := DB.QueryRow(`SELECT * FROM videos WHERE 
  course_id = ? ORDER BY id ASC LIMIT 1`, courseID)
	err := row.Scan(&video.ID, &video.Title, &video.Description,
		&video.VideoHLS, &video.Thumbnail, &video.Length, &video.Duration,
		&video.Views, &video.CourseID, &video.CreatedAt, &video.SReview)
	if err != nil {
		if err == sql.ErrNoRows {
			return video, fmt.Errorf("No video found with course id %s", courseID)
		}
		return video, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return video, nil
}

func UpdateVideo(v Video) error {
	result, err := DB.Exec(`UPDATE videos SET 
  title = ?, description = ?, video_hls = ?, thumbnail = ?, length = ?, duration = ?
  WHERE id = ?`,
		v.Title, v.Description, v.VideoHLS, v.Thumbnail, v.Length, v.Duration, v.ID)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No video found with id %d", v.ID)
	}
	return nil
}

func DeleteVideoByID(id string) error {
	result, err := DB.Exec("DELETE FROM videos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No video found with the id %s", id)
	}
	return nil
}

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
			AND (
		    LOWER(v.title) LIKE '%' || LOWER(?) || '%' 
			  OR LOWER(v.description) LIKE '%' || LOWER(?) || '%'
      )
		ORDER BY
			v.id;
	`

	rows, err := DB.Query(query, userId, userId, courseId, searchParam, searchParam)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
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
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return videos, nil
}

func GetAdminVideos(courseId, searchParam string) ([]Video, error) {
	var videos []Video
	rows, err := DB.Query(`SELECT * FROM videos
		WHERE course_id = ? AND (title LIKE ? OR description LIKE ?)
		ORDER BY id DESC
		`,
		courseId, searchParam, searchParam)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.VideoHLS,
			&v.Thumbnail, &v.Duration, &v.Length, &v.Views, &v.CourseID,
			&v.CreatedAt, &v.SReview); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		videos = append(videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return videos, nil
}

func CreateVideo(v Video) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO videos
  (title, description, video_hls, thumbnail, length, duration, course_id, created_at) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		v.Title, v.Description, v.VideoHLS, v.Thumbnail, v.Length, v.Duration, v.CourseID, date)

	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
