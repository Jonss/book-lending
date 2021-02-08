package pg

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Jonss/book-lending/infra/logger"
	_ "github.com/lib/pq"
)

func GetDbClient() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	client, err := sql.Open("postgres", connectionURL)
	if err != nil {
		logger.Error("Error on creating db connection" + err.Error())
		panic(err)
	}
	client.SetConnMaxIdleTime(time.Minute * 3)
	client.SetMaxIdleConns(10)
	client.SetMaxIdleConns(10)

	return client
}
