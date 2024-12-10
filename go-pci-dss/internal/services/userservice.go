package services

import (
	"database/sql"
	"errors"
	"fmt"
	"go-pci-dss/internal/models"
	"go-pci-dss/utils"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) RegisterUser(user models.User) (string, error) {

	if user.Username == "" || user.Password == "" {
		return "", errors.New("username and password are required")
	}
	// Hashovanje šifre
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GoApp",
		AccountName: fmt.Sprintf("%s@goapp.com", user.Username),
	})
	if err != nil {
		return "", err
	}

	// Unos korisnika u bazu
	_, err = s.DB.Exec(
		"INSERT INTO users (username, password, role, totpsecret) VALUES ($1, $2, $3, $4)",
		user.Username, hashedPassword, user.Role, key.Secret(),
	)
	if err != nil {
		return "", err
	}

	return key.URL(), nil
}

// Login - funkcija za autentifikaciju korisnika
func (s *UserService) Login(username, password string, otp string) (string, error) {
	// 1. Proveravamo da li postoji korisnik sa tim korisničkim imenom
	var user models.User
	err := s.DB.QueryRow("SELECT id, username, password, role, totpsecret FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.TOTPSecret)
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

	// 3. Proveravamo da li je TOTP kod tačan
	valid, err := utils.ValidateTOTP(user.TOTPSecret, otp)
	if err != nil || !valid {
		return "", errors.New("invalid or expired TOTP code")
	}

	// 4. Generišemo JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	// 5. Vraćamo token
	return token, nil
}
