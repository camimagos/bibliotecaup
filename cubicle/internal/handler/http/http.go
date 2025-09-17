package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"bibliotecaup.com/cubicle/internal/controller/cubicle"
)

type Handler struct {
	ctrl *cubicle.Controller
}

func New(ctrl *cubicle.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetCubicleDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, cubicle.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
