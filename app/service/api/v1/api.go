package v1

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/guarilha/go-ddd-starter/domain/user"
)

type ApiHandlers struct {
	UsersUseCase user.UseCase
}

func (a *ApiHandlers) Routes(router *chi.Mux) {
	router.Route("/api/v1", func(r chi.Router) {
		config := huma.DefaultConfig("Guaca API", "0.0.1")
		config.Servers = []*huma.Server{
			{
				URL: "/api/v1",
			},
		}

		api := humachi.New(r, config)

		huma.Post(api, "/signup", SignUpHandler(a.UsersUseCase))

	})
}
