package user

import (
	"github.com/guarilha/go-ddd-starter/domain"
	"github.com/urfave/cli/v3"
)

func Command(domains *domain.Domains) *cli.Command {
	newUser := NewUserCommand{}

	return &cli.Command{
		Name:  "user",
		Usage: "Manage and create users",
		Commands: []*cli.Command{
			newUser.Command(domains),
		},
	}
}
