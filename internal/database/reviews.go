package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type Review struct {
	ID          int64   `json:"id"`
	CourseId    int64   `json:"course_id"`
	UserId      int64   `json:"user_id"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Public      bool    `json:"public"`
	CourseTitle string  `json:"course_title"`
	CreatedAt   string  `json:"created_at"`
}

func DeleteReview(id string) error {
	result, err := DB.Exec(`DELETE FROM reviews WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No review found with the id %v", id)
	}
	return nil

}

func UpdatePublicStatus(public, id string) error {
	result, err := DB.Exec(`UPDATE reviews SET public = ? WHERE id = ?`, public, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No review found with the id %v", id)
	}
	return nil
}

func GetAdminReviews(searchParam string, isPublic string) ([]Review, error) {
	var reviews []Review
	query := `
		SELECT reviews.*, courses.title AS course_title
		FROM reviews
		JOIN courses ON reviews.course_id = courses.id
		WHERE 
			(LOWER(reviews.description) LIKE '%' || LOWER(?) || '%' 
			OR LOWER(reviews.author) LIKE '%' || LOWER(?) || '%'
			OR LOWER(courses.title) LIKE '%' || LOWER(?) || '%')
	`

	if isPublic != "" {
		if isPublic == "true" {
			query += " AND reviews.public = 'true'"
		} else {
			query += " AND reviews.public = 'false'"
		}
	}

	rows, err := DB.Query(query, searchParam, searchParam, searchParam)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.ID, &r.CourseId, &r.UserId, &r.Author, &r.Description, &r.Rating, &r.Public, &r.CreatedAt, &r.CourseTitle); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return reviews, nil
}

func GetPublicReviewsByCourseId(courseId string) ([]Review, error) {
	var reviews []Review
	rows, err := DB.Query(`SELECT * FROM reviews 
  WHERE course_id = ? AND public = 'true'`, courseId)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.ID, &r.CourseId, &r.UserId, &r.Author, &r.Description, &r.Rating, &r.Public, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return reviews, nil
}

func UserReviewExists(userId int64, courseId string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM reviews WHERE course_id = ? AND user_id = ?)"
	err := DB.QueryRow(query, courseId, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return exists, nil
}

func CreateReview(userId int64, courseId, author, description string, rating float64) (int64, error) {
	date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO reviews
  (user_id, course_id, author, description, rating, created_at) 
  VALUES (?, ?, ?, ?, ?, ?)`, userId, courseId, author, description, rating, date)
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return id, nil
}
