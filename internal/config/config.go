package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Env err", err)
		return err
	}
	return nil
}
