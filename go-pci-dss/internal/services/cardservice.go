package services

import (
	"database/sql"
	"errors"
	"go-pci-dss/internal/models"
	"go-pci-dss/utils"
	"log"
	"os"
	"regexp"
)

type CardholderService struct {
	DB *sql.DB
}

func NewCardholderService(db *sql.DB) *CardholderService {
	return &CardholderService{DB: db}
}

func (s *CardholderService) GetAllCardholders() ([]models.Cardholder, error) {

	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	rows, err := s.DB.Query(
		`SELECT id, name, 
			pgp_sym_decrypt(card_number::bytea, $1) AS card_number, 
			expiration_date 
		 FROM cardholders`, encryptionKey,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cardholders []models.Cardholder
	for rows.Next() {
		var c models.Cardholder
		if err := rows.Scan(&c.ID, &c.Name, &c.CardNumber, &c.ExpirationDate); err != nil {
			return nil, err
		}
		err1 := utils.InitializeKey()
		if err1 != nil {
			log.Fatalf("Failed to initialize key: %v", err1)
		}
		decryptedCardNumber := utils.Decrypt(c.CardNumber)

		if len(decryptedCardNumber) >= 4 {
			maskedCardNumber := "**** **** **** " + decryptedCardNumber[len(decryptedCardNumber)-4:]
			c.CardNumber = maskedCardNumber
		} else {
			c.CardNumber = "Invalid Card Number"
		}

		cardholders = append(cardholders, c)
	}
	return cardholders, nil
}

func (s *CardholderService) CreateCardholder(cardholder models.Cardholder) error {
	if cardholder.Name == "" || cardholder.CardNumber == "" || cardholder.ExpirationDate == "" || cardholder.CVV == "" {
		return errors.New("name, card number, expiration date, and CVV are required")
	}

	err1 := utils.InitializeKey()
	if err1 != nil {
		log.Fatalf("Failed to initialize key: %v", err1)
	}

	encryptedCVV := utils.Encrypt(cardholder.CVV)
	encryptedPAN := utils.Encrypt(cardholder.CardNumber)

	if err := s.AuthorizeCard(cardholder.CardNumber, encryptedCVV); err != nil {
		return err
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	_, err := s.DB.Exec(
		`INSERT INTO cardholders (name, card_number, expiration_date) 
		 VALUES ($1, pgp_sym_encrypt($2, $3), $4)`,
		cardholder.Name, encryptedPAN, encryptionKey, cardholder.ExpirationDate,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *CardholderService) AuthorizeCard(cardNumber, encryptedCVV string) error {
	matched, err := regexp.MatchString(`^\d{16}$`, cardNumber)
	if err != nil {
		return errors.New("error while validating card number")
	}
	if !matched {
		return errors.New("invalid card number: must be exactly 16 digits")
	}

	cvv := utils.Decrypt(encryptedCVV)

	matchedCVV, err := regexp.MatchString(`^\d{3,4}$`, cvv)
	if err != nil {
		return errors.New("error while validating CVV")
	}
	if !matchedCVV {
		return errors.New("invalid CVV: must be 3 or 4 digits")
	}
	return nil
}
