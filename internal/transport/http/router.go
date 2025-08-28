package transport

import (
	"net/http"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/transport/http/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(repo repository.Repository, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	h := handlers.NewDeviceHandler(repo, cfg)
	r.Get("/api/devices", h.GetDevices)

	// health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("ok")) })
	return r
}
