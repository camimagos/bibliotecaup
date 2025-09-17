package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"bibliotecaup.com/metadata/internal/controller/metadata"
	repository "bibliotecaup.com/metadata/internal/repository"
)

type Handler struct {
	controller *metadata.Controller
}

func New(controller *metadata.Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest) // devolver un 400
	}

	ctx := r.Context()
	m, err := h.controller.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound) // devolver un 404
		return
	} else if err != nil {
		log.Printf("failed to get metadata. repository error: %v", err)
		w.WriteHeader(http.StatusInternalServerError) // devolver un 500
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("failed to encode metadata response: %v", err)
	}
}
