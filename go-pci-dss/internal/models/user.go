package models

type User struct {
	ID         int    `json:"id"`       // Jedinstveni identifikator korisnika
	Username   string `json:"username"` // Korisničko ime za login
	Password   string `json:"password"` // Šifrovana lozinka
	Role       string `json:"role"`     // Uloga korisnika (npr. "admin", "user")
	TOTPSecret string `json:"totpsecret"`
}
