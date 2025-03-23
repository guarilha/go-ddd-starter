package api

import (
	"net/http"
	"text/template"

	"github.com/guarilha/go-ddd-starter/domain"
	"github.com/guarilha/go-ddd-starter/domain/user"
)

func UserGet(d *domain.Domains) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		email := r.URL.Query().Get("email")

		if name == "" || email == "" {
			http.Error(w, "Name and email are required", http.StatusBadRequest)
			return
		}

		user, err := d.User.SignUp(r.Context(), user.SignUpParams{
			Name:  name,
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templates, "templates/user_view.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
