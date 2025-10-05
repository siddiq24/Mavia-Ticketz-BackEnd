package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error, string) {
	dbUser := os.Getenv("DB_USER")
	// dbUser := "postgres"
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	// dbHost := "localhost"
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format: postgres://username:password@host:port/dbname
	cons := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := pgxpool.New(context.Background(), cons)
	if err != nil {
		return nil, err, "Not Found"
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err, "Not Found"
	}

	return db, nil, dbName
}
