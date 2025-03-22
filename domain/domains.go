package domain

import (
	"fmt"

	"github.com/guarilha/go-ddd-starter/domain/user"
	"github.com/guarilha/go-ddd-starter/gateways/pg"
)

type Config struct {
	Example string
}

type Domains struct {
	User *user.Domain
}

func NewDomains(db pg.Db, cfg Config) (*Domains, error) {

	userDomain, err := user.New(db, user.Config{
		Example: cfg.Example,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user domain: %v", err)
	}

	return &Domains{
		User: userDomain,
	}, nil
}
