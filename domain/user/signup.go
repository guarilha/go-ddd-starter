package user

import (
	"context"

	"github.com/guarilha/go-ddd-starter/domain/entities"
)

func (uc *UseCase) SignUp(ctx context.Context, name, email string) (*entities.User, error) {
	user, err := entities.NewUser(name, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
