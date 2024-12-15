package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatalf("missing or empty env DATABASE_URL")
	}

	db, err := sql.Open("pgx", dbUrl)

	if err != nil {
		log.Fatalf("failed to connect to database, Error -> %v\n", err.Error())
	}

	return db, db.Ping()
}
