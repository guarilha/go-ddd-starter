package api

import (
	"net/http"
	"sync"
	"text/template"
)

type CounterHandler struct {
	counter   int
	mu        sync.Mutex
	templates *template.Template
}

func NewCounterHandler(templates *template.Template) *CounterHandler {
	return &CounterHandler{
		counter:   0,
		templates: templates,
	}
}

func (h *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	action := r.URL.Query().Get("action")
	switch action {
	case "increment":
		h.counter++
	case "decrement":
		h.counter--
	}

	h.templates.ExecuteTemplate(w, "counter.html", h.counter)
}
