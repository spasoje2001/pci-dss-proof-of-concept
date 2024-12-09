package services

import (
	"database/sql"
	"errors"
	"go-pci-dss/internal/models"
	"go-pci-dss/utils"
	"regexp"
)

type CardholderService struct {
	DB *sql.DB
}

func NewCardholderService(db *sql.DB) *CardholderService {
	return &CardholderService{DB: db}
}

func (s *CardholderService) GetAllCardholders() ([]models.Cardholder, error) {
	rows, err := s.DB.Query("SELECT id, name, card_number, expiration_date FROM cardholders")
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

		// Desifrujemo broj kartice
		decryptedCardNumber := utils.Decrypt(c.CardNumber)

		// Maskiramo broj kartice
		if len(decryptedCardNumber) >= 4 {
			maskedCardNumber := "**** **** **** " + decryptedCardNumber[len(decryptedCardNumber)-4:]
			c.CardNumber = maskedCardNumber
		} else {
			// Ako je broj kartice iz nekog razloga prekratak, samo postavimo placeholder
			c.CardNumber = "Invalid Card Number"
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

	encryptedCVV := utils.Encrypt(cardholder.CVV)        // Ovde koristimo šifrovanje
	encryptedPAN := utils.Encrypt(cardholder.CardNumber) // Ovde koristimo šifrovanje

	// Simuliramo autorizaciju koristeći šifrovani CVV
	if err := s.AuthorizeCard(cardholder.CardNumber, encryptedCVV); err != nil {
		return err
	}

	// Pravimo SQL upit za unos podataka u bazu
	_, err := s.DB.Exec(
		"INSERT INTO cardholders (name, card_number, expiration_date) VALUES ($1, $2, $3)",
		cardholder.Name, encryptedPAN, cardholder.ExpirationDate,
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

	// Desifrovanje CVV
	cvv := utils.Decrypt(encryptedCVV)

	// Provera da li CVV ima 3 ili 4 cifre
	matchedCVV, err := regexp.MatchString(`^\d{3,4}$`, cvv)
	if err != nil {
		return errors.New("error while validating CVV")
	}
	if !matchedCVV {
		return errors.New("invalid CVV: must be 3 or 4 digits")
	}
	return nil
}
