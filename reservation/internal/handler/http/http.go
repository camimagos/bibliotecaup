package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"bibliotecaup.com/reservation/internal/controller/reservation"
	"bibliotecaup.com/reservation/pkg/model"
)

type Handler struct {
	ctrl *reservation.Controller
}

func New(ctrl *reservation.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))

	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPost: // Crear una reservación
		if err := req.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		res := &model.Reservation{
			RecordID:   model.RecordID(req.Form.Get("recordId")),
			RecordType: model.RecordType(req.Form.Get("recordType")),
			UserID:     model.UserID(req.Form.Get("userId")),
			Start:      parseTime(req.Form.Get("start")),
			End:        parseTime(req.Form.Get("end")),
			Status:     model.Status(req.Form.Get("status")),
		}

		if err := h.ctrl.PutReservation(req.Context(), res.RecordID, res.RecordType, res); err != nil {
			log.Printf("Error creating reservation: %v", err)
			if errors.Is(err, reservation.ErrInvalidData) {
				http.Error(w, "Invalid reservation data", http.StatusBadRequest)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodGet: // Disponibilidad próxima
		recordID := model.RecordID(req.URL.Query().Get("recordId"))
		recordType := model.RecordType(req.URL.Query().Get("recordType"))

		availability, err := h.ctrl.Reservation(req.Context(), recordID, recordType)
		if err != nil {
			log.Printf("Error fetching availability: %v", err)
			if errors.Is(err, reservation.ErrNotFound) {
				http.Error(w, "Record not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(availability); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// helper para parsear fechas ISO
func parseTime(s string) (t time.Time) {
	t, _ = time.Parse(time.RFC3339, s)
	return
}
