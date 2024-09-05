package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func CreateReview(input database.Review) (database.Review, error) {
	if input.Rating == 0 {
		return database.Review{}, fmt.Errorf("The rating is required.")
	}

	if input.Description == "" {
		return database.Review{}, fmt.Errorf("The name is required.")
	}
	if len(input.Description) > 255 {
		return database.Review{}, fmt.Errorf("The description should not have more than 255 characters.")
	}

	return database.Review{
		Rating:      input.Rating,
		Description: input.Description,
	}, nil
}
