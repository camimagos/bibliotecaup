package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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
	switch req.Method {
	case http.MethodPost: // Crear una reservación
		var res model.Reservation
		if err := json.NewDecoder(req.Body).Decode(&res); err != nil {
			log.Printf("Error decoding reservation: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if res.RecordID == "" || res.RecordType == "" || res.UserID == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutReservation(req.Context(), res.RecordID, res.RecordType, &res); err != nil {
			log.Printf("Error creating reservation: %v", err)
			if errors.Is(err, reservation.ErrInvalidData) {
				http.Error(w, "Invalid reservation data", http.StatusBadRequest)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		//log.Printf("Reservation created successfully: %+v", res)
		w.WriteHeader(http.StatusCreated)

	case http.MethodGet: // Disponibilidad próxima
		recordID := model.RecordID(req.URL.Query().Get("recordId"))
		recordType := model.RecordType(req.URL.Query().Get("recordType"))

		if recordID == "" || recordType == "" {
			log.Printf("Missing required parameters: recordID=%s, recordType=%s", recordID, recordType)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
// func parseTime(s string) (t time.Time) {
// 	t, _ = time.Parse(time.RFC3339, s)
// 	return
// }
