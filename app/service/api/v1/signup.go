package v1

import (
	"context"

	"github.com/guarilha/go-ddd-starter/domain/user"
)

type SignUpInput struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"User email is required."`
		Name  string `json:"name" example:"John Doe" doc:"User name is required."`
	}
}

type SignUpOutput struct {
	Body user.User
}

func SignUpHandler(uc *user.Domain) func(ctx context.Context, input *SignUpInput) (*SignUpOutput, error) {
	return func(ctx context.Context, input *SignUpInput) (*SignUpOutput, error) {
		user, err := uc.SignUp(ctx, user.SignUpParams{
			Email: input.Body.Email,
			Name:  input.Body.Name,
		})
		if err != nil {
			return nil, err
		}

		output := &SignUpOutput{
			Body: *user,
		}

		return output, nil
	}
}
