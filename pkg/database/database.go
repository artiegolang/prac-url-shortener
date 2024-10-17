package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	log.Println("Connected to database")
	return &DB{Pool: pool}
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
