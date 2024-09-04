package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type File struct {
	ID        int64  `json:"id"`
	Path      string `json:"path"`
	VideoID   string `json:"video_id"`
	Page      string `json:"page"`
	CreatedAt string `json:"created_at"`
}

func FileExistsByVideoId(videoId int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM files WHERE video_id = ?)"
	err := DB.QueryRow(query, videoId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteFileByID(id string) error {
	result, err := DB.Exec("DELETE FROM files WHERE id = ?", id)
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

func GetTotalPagesByVideoId(video_id string) (int, error) {
	var total int
	row := DB.QueryRow(`SELECT COUNT(DISTINCT page) AS total_pages 
  FROM files WHERE video_id = ?;`, video_id)
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return total, nil
}

func GetFiles(videoID, page string) ([]File, error) {
	var files []File
	rows, err := DB.Query(`SELECT * FROM files
  WHERE video_id = ? AND page = ? 
  ORDER BY id`,
		videoID, page)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.Path, &f.Page, &f.VideoID, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return files, nil
}

func CreateFile(f File) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`INSERT INTO files
   (path, video_id, page, created_at) 
    VALUES (?, ?, ?, ?)`,
		f.Path, f.VideoID, f.Page, date)

	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
