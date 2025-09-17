package memory

import (
	"context"
	"fmt"

	model "bibliotecaup.com/metadata/pkg"
)

// Repository es una implementación en memoria del repositorio de metadatos.
type Repository struct {
	data map[string]*model.Metadata
}

// New crea una nueva instancia del repositorio en memoria.
func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

// Implementa los métodos necesarios para la interfaz del repositorio.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	if metadata, ok := r.data[id]; ok {
		return metadata, nil
	}
	return nil, fmt.Errorf("metadata not found")
}
