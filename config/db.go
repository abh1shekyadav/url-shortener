package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("Database url is required")
	}
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Println("Unable to connect to database: ", err)
	}
	DB = pool
	fmt.Println("Connected to PostgreSQL")
}
