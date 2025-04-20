package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() error {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Println("Env err", err)
		return err
	}
	return nil
}
