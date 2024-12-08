package models

type Cardholder struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
