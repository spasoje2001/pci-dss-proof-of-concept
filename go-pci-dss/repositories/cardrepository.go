package repositories

import (
	"errors"
	"go-pci-dss/models"
)

type CardRepository interface {
	Save(card models.Card) error
	GetByID(id string) (*models.Card, error)
}

type MemoryCardRepository struct {
	data map[string]models.Card
}

func NewCardRepository() *MemoryCardRepository {
	return &MemoryCardRepository{data: make(map[string]models.Card)}
}

func (r *MemoryCardRepository) Save(card models.Card) error {
	r.data[card.ID] = card
	return nil
}

func (r *MemoryCardRepository) GetByID(id string) (*models.Card, error) {
	card, exists := r.data[id]
	if !exists {
		return nil, errors.New("card not found")
	}
	return &card, nil
}
