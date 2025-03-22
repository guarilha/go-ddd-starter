package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/guarilha/go-ddd-starter/domain/user"
)

type ApiHandlers struct {
	UserDomain *user.Domain
}

func (a *ApiHandlers) Routes(router *chi.Mux) {
	router.Route("/api/v1", func(r chi.Router) {
		config := huma.DefaultConfig("Project API", "0.0.1")
		config.Servers = []*huma.Server{
			{
				URL: "/api/v1",
			},
		}

		api := humachi.New(r, config)

		huma.Register(api, huma.Operation{
			Method:      http.MethodPost,
			Path:        "/signup",
			Summary:     "Sign up",
			Tags:        []string{"user"},
			Description: "Sign up a new user",
		}, SignUpHandler(a.UserDomain))

	})
}
