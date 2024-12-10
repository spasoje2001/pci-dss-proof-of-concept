package config

import (
	"log"
	"os"

	"github.com/casbin/casbin"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
	CasbinEnforcer     *casbin.Enforcer // Casbin enforcer
}

func LoadConfig() (*Config, error) {
	// Učitaj .env fajl
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Inicijalizuj Casbin enforcer
	enforcer, err := initCasbin()
	if err != nil {
		return nil, err
	}

	// Vraćanje konfiguracije sa Casbin enforcerom
	return &Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		CasbinEnforcer:     enforcer,
	}, nil
}

// initCasbin - funkcija koja inicijalizuje Casbin enforcer
func initCasbin() (*casbin.Enforcer, error) {
	// Putanja do modela i politike
	modelPath := "casbin_model.conf"
	policyPath := "policy.csv"

	// Inicijalizacija Casbin enforcera sa modelom i politikom
	enforcer := casbin.NewEnforcer(modelPath, policyPath)

	return enforcer, nil
}
