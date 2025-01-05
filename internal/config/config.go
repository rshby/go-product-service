package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func GetEnv(k string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	return os.Getenv(k)
}

func Port() int {
	if p := GetEnv("PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			return port
		}
	}

	return 4002
}
