package model

import (
	model "bibliotecaup.com/metadata/pkg"
	reservationModel "bibliotecaup.com/reservation/pkg/model"
)

type CubicleDetails struct {
	Metadata    model.Metadata                 `json:"metadata"`
	Reservation *reservationModel.Availability `json:"reservation,omitempty"`
}
