package main

import (
	"fmt"
	"os"

	"github.com/Jonss/book-lending/app"
	"github.com/Jonss/book-lending/infra/logger"
	_ "github.com/lib/pq"
)

func main() {
	logger.Info(fmt.Sprintf("Starting application in environment: [%s]", os.Getenv("ENV")))

	app.Start()
	fmt.Println("Ahoy World")
}
