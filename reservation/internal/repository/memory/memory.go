package memory

import (
	"context"
	"log"

	"bibliotecaup.com/reservation/internal/repository"
	"bibliotecaup.com/reservation/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Reservation
}

func New() *Repository {
	return &Repository{
		//map[model.RecordType]map[model.RecordID][]model.Reservation{},
		data: make(map[model.RecordType]map[model.RecordID][]model.Reservation),
	}
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Reservation, error) {
	recordTypeData, ok := r.data[recordType]
	if !ok {
		return nil, repository.ErrNotFound
	}
	reservations, ok := recordTypeData[recordID]
	if !ok || len(reservations) == 0 {
		return nil, repository.ErrNotFound
	}
	return reservations, nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, reservation *model.Reservation) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Reservation{}
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *reservation)
	log.Printf("Reservation saved: %+v", reservation)
	return nil
}
