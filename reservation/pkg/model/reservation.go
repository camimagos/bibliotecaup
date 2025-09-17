package model

import "time"

type RecordID string

type RecordType string

type UserID string
type Status string

const (
	RecordTypeCubicle RecordType = "cubicle"
	RecordTypeCabin   RecordType = "cabin"

	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

// Reservation representa una reservación en el sistema
type Reservation struct {
	RecordID   RecordID   `json:"recordId"`
	RecordType RecordType `json:"recordType"`
	UserID     UserID     `json:"userId"`
	Start      time.Time  `json:"start"`
	End        time.Time  `json:"end"`
	Status     Status     `json:"status"`
}

type Availability struct {
	AvailableNow  bool      `json:"availableNow"`  //si está disponible en este momento
	NextAvailable time.Time `json:"nextAvailable"` //próxima hora disponible
}
