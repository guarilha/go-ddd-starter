package user

import (
	"fmt"

	"github.com/guarilha/go-ddd-starter/domain/user/repository"
	"github.com/guarilha/go-ddd-starter/gateways/pg"
)

type Domain struct {
	R repository.Querier
}

type Config struct {
	Example string
}

func New(db pg.Db, cfg Config) (*Domain, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}

	r := repository.New(db)

	return &Domain{
		R: r,
	}, nil
}
