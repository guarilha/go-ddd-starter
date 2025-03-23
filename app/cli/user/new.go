package user

import (
	"context"
	"fmt"

	"github.com/guarilha/go-ddd-starter/domain"
	"github.com/guarilha/go-ddd-starter/domain/user"
	"github.com/urfave/cli/v3"
)

type NewUserCommand struct{}

func (c NewUserCommand) Command(domains *domain.Domains) *cli.Command {
	return &cli.Command{
		Name:        "new",
		Aliases:     []string{"create"},
		Description: "Create a new user and store the credentials locally. By default we generate a new one and store it locally to use with the new User.",
		UsageText:   "cli user new",
		Action:      c.Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "Test User",
				Usage: "User name",
			},
			&cli.StringFlag{
				Name:  "email",
				Usage: "User email",
			},
		},
	}
}

func (c NewUserCommand) Action(ctx context.Context, cmd *cli.Command) error {

	domains := ctx.Value("domains").(*domain.Domains)

	user, err := domains.User.SignUp(ctx, user.SignUpParams{
		Name:  cmd.String("name"),
		Email: cmd.String("email"),
	})
	if err != nil {
		return err
	}

	fmt.Println("User created successfully")
	fmt.Println("Name:", user.Name)
	fmt.Println("Email:", user.Email)
	fmt.Println("ID:", user.ID)

	return nil
}
