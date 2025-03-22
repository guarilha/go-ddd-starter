package user

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type SignUpParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (d *Domain) SignUp(ctx context.Context, params SignUpParams) (*User, error) {
	user, err := NewUser(params.Name, params.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
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
