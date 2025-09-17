package cubicle

import (
	"context"
	"errors"
	"log"

	"bibliotecaup.com/cubicle/internal/gateway"
	"bibliotecaup.com/cubicle/pkg/model"
	metadataModel "bibliotecaup.com/metadata/pkg"
	reservationModel "bibliotecaup.com/reservation/pkg/model"
)

var ErrNotFound = errors.New("cubicle: not found")

type reservationGateway interface {
	GetAggregatedReservation(ctx context.Context, recordID reservationModel.RecordID, recordType reservationModel.RecordType) (*reservationModel.Availability, error)
	PutReservation(ctx context.Context, recordID reservationModel.RecordID, recordType reservationModel.RecordType, reservation *reservationModel.Reservation) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	reservationGateway reservationGateway
	metadataGateway    metadataGateway
}

func New(reservationGateway reservationGateway, metametadataGateway metadataGateway) *Controller {
	return &Controller{reservationGateway, metametadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.CubicleDetails, error) {
	// Obtener los metadatos del cubículo
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Inicializar los detalles del cubículo con los metadatos
	details := &model.CubicleDetails{Metadata: *metadata}

	// Obtener la disponibilidad agregada de las reservaciones
	availability, err := c.reservationGateway.GetAggregatedReservation(ctx, reservationModel.RecordID(id), reservationModel.RecordTypeCubicle)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		log.Printf("No availability found for cubicle ID %s", id)
	} else if err != nil {
		return nil, err
	} else {
		details.Reservation = availability
	}

	return details, nil
}

// type Controller struct {
//     metadataGateway    *http.Gateway
//     reservationGateway *reservationGateway.Gateway
// }

// func New(metadataGW *http.Gateway, reservationGW *reservationGateway.Gateway) *Controller {
//     return &Controller{
//         metadataGateway:    metadataGW,
//         reservationGateway: reservationGW,
//     }
// }

// type Aggregated struct {
//     Reservations *reservationGateway.Availability `json:"reservations"`
//     Metadata     *model.Metadata                  `json:"metadata"`
// }

// func (c *Controller) Get(ctx context.Context, id string) (*Aggregated, error) {
//     metadata, err := c.metadataGateway.Get(ctx, id)
//     if err != nil {
//         return nil, err
//     }

//     availability, err := c.reservationGateway.GetAggregated(ctx, id)
//     if err != nil {
//         return nil, err
//     }

//     return &Aggregated{
//         Reservations: availability,
//         Metadata:     metadata,
//     }, nil
// }
