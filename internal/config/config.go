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

func initConfig() Config {
	var DB_NAME, uri string = "dairyDB", "mongodb://localhost:27017/"
	if err := godotenv.Load(); err != nil {
		log.Println("Set your 'MONGODB_URI' environment variable. " + "No .env file found\nUsing the default 'mongodb://localhost:27017/'")
	}

	conf := Config{
		PORT:       getEnv("PORT", "3000"),
		DB_URI:     getEnv("MONGODB_URI", uri),
		DB_NAME:    getEnv("DB_NAME", DB_NAME),
		JWT_SECRET: getEnv("JWT_SECRET", "randomjwtsecretkey"),
		ENV:        getEnv("ENV", "dev"),
	}
	return conf
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
