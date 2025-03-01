package entities

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(name, email string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
}
