package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironment() error {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("failed to load env, err %v", err)
	}
	return nil
}
