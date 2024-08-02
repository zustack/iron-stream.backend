package database

import "fmt"

type File struct {
  ID          int64  `json:"id"`
  Path        string `json:"path"`
  VideoID     int64  `json:"video_id"` 
  Page int64  `json:"page"`
  SortOrder   int64  `json:"sort_order"`
  CreatedAt   string `json:"created_at"`
}

func DeleteFileByID(id int64) error {
	result, err := DB.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
    return fmt.Errorf("DeleteFileByID: video id: %d: %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
    return fmt.Errorf("DeleteFileByID: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
    return fmt.Errorf("DeleteFileByID: no file found with ID: %d", id)
	}
	return nil
}

func GetFiles(video_id, page int64) ([]File, error) {
	var files []File
	rows, err := DB.Query(`SELECT * FROM files
  WHERE video_id = ? AND page = ? 
  ORDER BY sort_order`,
		video_id, page)
	if err != nil {
		return nil, fmt.Errorf("GetFiles: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.Path, &f.VideoID, &f.Page, &f.SortOrder, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetFiles: %v", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFiles: %v", err)
	}
	return files, nil
}

func CreateFile(f File) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO files
  (path, video_id, page, sort_order) 
  VALUES (?, ?, ?, ?)`,
  f.Path, f.VideoID, f.Page, 0)

	if err != nil {
		return 0, fmt.Errorf("CreateFile: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateFile: %v", err)
	}

	_, err = DB.Exec("UPDATE files SET sort_order = ? WHERE id = ?", id, id)
	if err != nil {
		return 0, fmt.Errorf("CreateFile/updateSortOrder: %v", err)
	}

	return id, nil
}
