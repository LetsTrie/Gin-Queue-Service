package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type config struct {
	Port    			string `validate:"required"`
	RedisHost     string `validate:"required"`
	RedisPassword string `validate:"required"`
}

var Env config

func Init() {
	err := godotenv.Load()

	if(err != nil) {
		log.Fatal("Error loading .env file!")
	}

	Env = config{
		Port:					 os.Getenv("PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}

	validate := validator.New()
	if err := validate.Struct(Env); err != nil {
		log.Fatalf("❌ Configuration validation failed: %s\n", err)
	}
}