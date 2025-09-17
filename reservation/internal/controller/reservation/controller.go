package reservation

import (
	"context"
	"errors"
	"sort"
	"time"

	"bibliotecaup.com/reservation/pkg/model"
)

var ErrNotFound = errors.New("reservation: not found")
var ErrInvalidData = errors.New("invalid data")

type reservationRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Reservation, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, reservation *model.Reservation) error
}

type Controller struct {
	repo reservationRepository
}

func New(repo reservationRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Reservation(ctx context.Context, id model.RecordID, t model.RecordType) (*model.Availability, error) {
	reservations, err := c.repo.Get(ctx, id, t)
	if errors.Is(err, ErrNotFound) {
		return &model.Availability{
			AvailableNow:  true,
			NextAvailable: time.Now(),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// ordenar por hora de inicio
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].Start.Before(reservations[j].Start)
	})

	// ver si está ocupado ahora
	for _, res := range reservations {
		if res.Start.Before(now) && res.End.After(now) {
			return &model.Availability{
				AvailableNow:  false,
				NextAvailable: res.End,
			}, nil
		}
	}

	// buscar la próxima reservación futura
	for _, res := range reservations {
		if res.Start.After(now) {
			return &model.Availability{
				AvailableNow:  true,
				NextAvailable: res.Start,
			}, nil
		}
	}

	// no hay reservaciones futuras
	return &model.Availability{
		AvailableNow:  true,
		NextAvailable: now,
	}, nil
}

func (c *Controller) PutReservation(ctx context.Context, recordID model.RecordID, recordType model.RecordType, reservation *model.Reservation) error {
	return c.repo.Put(ctx, recordID, recordType, reservation)
}
