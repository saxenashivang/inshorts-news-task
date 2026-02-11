package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivangsaxena/inshorts-task/config"

	"github.com/shivangsaxena/inshorts-task/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.Init()

	// Database Connection (pgx)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Enable PostGIS extension
	if _, err := dbPool.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS postgis"); err != nil {
		log.Printf("Warning: Failed to enable PostGIS extension: %v", err)
	}
}
