package models

type Card struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	Name   string `json:"name"`
	Expiry string `json:"expiry"`
}
