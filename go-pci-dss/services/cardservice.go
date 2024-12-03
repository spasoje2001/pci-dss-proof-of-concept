package services

import (
	"errors"
	"go-pci-dss/models"
	"go-pci-dss/repositories"
	"go-pci-dss/utils"
)

type CardService struct {
	repo repositories.CardRepository
}

func NewCardService(repo repositories.CardRepository) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) SaveCard(card models.Card) error {
	if card.Number == "" {
		return errors.New("card number is required")
	}
	card.Number = utils.Encrypt(card.Number)
	return s.repo.Save(card)
}

func (s *CardService) GetCard(id string) (*models.Card, error) {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	card.Number = utils.Decrypt(card.Number)
	return card, nil
}
