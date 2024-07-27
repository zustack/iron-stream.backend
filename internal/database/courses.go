package database

import "fmt"

type Course struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Thumbnail   string `json:"thumbnail"`
	Preview     string `json:"preview"`
	Rating      int    `json:"rating"`
	NumReviews  int    `json:"num_reviews"`
	Duration    string `json:"duration"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

func CreateCourse(c Course) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO course
  (title, description, author, thumbnail, preview, duration, is_active, sort_order) 
  VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.Author, c.Thumbnail, c.Preview, c.Duration, c.IsActive, 0)

	if err != nil {
		return 0, fmt.Errorf("CreateCourse: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateCourse: %v", err)
	}

	_, err = DB.Exec("UPDATE course SET sort_order = ? WHERE id = ?", id, id)
	if err != nil {
		return 0, fmt.Errorf("CreateCourse/updateOrder: %v", err)
	}

	return id, nil
}
