package main

import (
	"attendance/internal/routers"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	connStr := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	db, err := GetDB(connStr)
	fmt.Println(err)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	router := routers.SetupRouter(db)

	router.Run("0.0.0.0:" + port)
}

func GetDB(uri string) (*sqlx.DB, error) {
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 30
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Test the connection
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("Failed to acquire connection from pool:", err)
		return nil, err
	}

	// To be safe
	defer conn.Release()

	// Get a connection
	db := stdlib.OpenDBFromPool(pool)
	if db == nil {
		return nil, fmt.Errorf("failed to open db")
	}

	return sqlx.NewDb(db, "postgres"), nil

}
