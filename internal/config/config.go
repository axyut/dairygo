package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT       string
	DB_URI     string
	DB_NAME    string
	JWT_SECRET string
	ENV        string
}

var Envs = initConfig()

// ld flags production
var uri string
var db_name string

func initConfig() Config {
	if uri == "" {
		uri = "mongodb://localhost:27017/"
	}

	if db_name == "" {
		db_name = "dairygo"
	}

	// .env file development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using the default or build values.")
	}

	conf := Config{
		PORT:       getEnv("PORT", "3000"),
		DB_URI:     getEnv("MONGODB_URI", uri),
		DB_NAME:    getEnv("DB_NAME", db_name),
		JWT_SECRET: getEnv("JWT_SECRET", "randomjwtsecretkey"),
		ENV:        getEnv("ENV", "dev"),
	}
	return conf
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		// fmt.Println(ok, value)
		return value
	}

	return fallback
}
