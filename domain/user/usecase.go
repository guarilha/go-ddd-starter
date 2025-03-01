package user

import "github.com/guarilha/go-ddd-starter/gateways/repository"

type UseCase struct {
	R repository.Querier
}
