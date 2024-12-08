package services

import (
	"database/sql"
	"errors"
	"go-pci-dss/internal/models"
)

type CardholderService struct {
	DB *sql.DB
}

func NewCardholderService(db *sql.DB) *CardholderService {
	return &CardholderService{DB: db}
}

func (s *CardholderService) GetAllCardholders() ([]models.Cardholder, error) {
	rows, err := s.DB.Query("SELECT id, name, card_number, expiration_date, cvv FROM cardholders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cardholders []models.Cardholder
	for rows.Next() {
		var c models.Cardholder
		if err := rows.Scan(&c.ID, &c.Name, &c.CardNumber, &c.ExpirationDate, &c.CVV); err != nil {
			return nil, err
		}
		cardholders = append(cardholders, c)
	}
	return cardholders, nil
}

func (s *CardholderService) CreateCardholder(cardholder models.Cardholder) error {
	// Proveravamo da li su ime i broj kartice prisutni
	if cardholder.Name == "" || cardholder.CardNumber == "" || cardholder.ExpirationDate == "" || cardholder.CVV == "" {
		return errors.New("name, card number, expiration date, and CVV are required")
	}

	// Pravimo SQL upit za unos podataka u bazu
	_, err := s.DB.Exec(
		"INSERT INTO cardholders (name, card_number, expiration_date, cvv) VALUES ($1, $2, $3, $4)",
		cardholder.Name, cardholder.CardNumber, cardholder.ExpirationDate, cardholder.CVV,
	)
	return err
}
