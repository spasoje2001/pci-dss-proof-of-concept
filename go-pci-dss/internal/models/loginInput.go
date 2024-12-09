package models

// LoginInput sadrži podatke koje korisnik unosi prilikom logovanja
type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	TOTPSecret string `json:"totpsecret"`
}
