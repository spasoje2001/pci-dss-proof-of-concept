package services

import (
	"database/sql"
	"errors"
	"go-pci-dss/internal/models"
	"go-pci-dss/utils"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) RegisterUser(user models.User) error {
	log.Printf("U servicu")
	// Validacija ulaza

	log.Printf(user.Username)
	log.Printf(user.Password)
	if user.Username == "" || user.Password == "" {
		return errors.New("username and password are required")
	}
	// Hashovanje šifre
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Unos korisnika u bazu
	_, err = s.DB.Exec(
		"INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		user.Username, hashedPassword, user.Role,
	)
	if err != nil {
		return err
	}

	return nil
}

// Login - funkcija za autentifikaciju korisnika
func (s *UserService) Login(username, password string) (string, error) {
	// 1. Proveravamo da li postoji korisnik sa tim korisničkim imenom
	var user models.User
	err := s.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// 2. Proveravamo da li je lozinka tačna
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 3. Generišemo JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	// 4. Vraćamo token
	return token, nil
}
