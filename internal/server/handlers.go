package server

import "net/http"

type Handlers struct {
	store string
}

func NewHandlers() *Handlers {
	return &Handlers{
		store: "",
	}
}

func (h *Handlers) Ping(w http.ResponseWriter) {

	w.WriteHeader(http.StatusOK)

}
