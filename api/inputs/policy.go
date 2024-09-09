package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func CreatePolicy(p database.Policy) (database.Policy, error) {
	if p.Content == "" {
		return database.Policy{}, fmt.Errorf("The policy content is required.")
	}
	if len(p.Content) > 55 {
		return database.Policy{}, fmt.Errorf("The policy content should not have more than 55 characters.")
	}
	if p.PType == "" {
		return database.Policy{}, fmt.Errorf("The policy type is required.")
	}
	if len(p.PType) > 455 {
		return database.Policy{}, fmt.Errorf("The policy type should not have more than 455 characters.")
	}
	return database.Policy{
    Content:     p.Content,
    PType:     p.PType,
	}, nil
}

