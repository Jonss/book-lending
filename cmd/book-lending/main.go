package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Jonss/book-lending/app"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	startEnvVar()
	logger.Info(fmt.Sprintf("Starting application in environment: [%s]", os.Getenv("ENV")))

	app.Start()
}

func startEnvVar() {
	err := godotenv.Load("env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
