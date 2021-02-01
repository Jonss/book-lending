package pg

import (
	"database/sql"

	"github.com/Jonss/book-lending/infra/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func MigratePSQL(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("Error creating postgres Driver on migration")
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://adapters/storages/pg/migrations", "postgres", driver)
	if err != nil {
		logger.Error("Error creating migration database instance")
		panic(err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			logger.Error("Error migrating")
			panic(err)
		}
	}
	logger.Info("Executing migration!")
}
