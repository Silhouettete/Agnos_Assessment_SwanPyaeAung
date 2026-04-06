package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("database is unreachable:", err)
	}

	log.Println("Connected to database")
	return pool
}
