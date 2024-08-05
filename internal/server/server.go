package server

import (
	"github.com/lionslon/yap-gophermart/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitServer(h *Handlers, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    cfg.Address,
		Handler: initRouter(h),
	}
}

func initRouter(h *Handlers) *chi.Mux {

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		h.Ping(w)
	})

	return router

}
