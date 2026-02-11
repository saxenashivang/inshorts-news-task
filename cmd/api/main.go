package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivangsaxena/inshorts-task/config"
	"github.com/shivangsaxena/inshorts-task/internal/adapter/storage/repository"
	"github.com/shivangsaxena/inshorts-task/internal/core/service"

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

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		source VARCHAR(255),
		url TEXT,
		category TEXT[],
		lat DOUBLE PRECISION,
		lng DOUBLE PRECISION,
		geom GEOGRAPHY(Point, 4326),
		published_at TIMESTAMP WITH TIME ZONE,
		relevance_score DOUBLE PRECISION,
		search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || content)) STORED
	);
	CREATE INDEX IF NOT EXISTS article_geom_idx ON articles USING GIST (geom);
	CREATE INDEX IF NOT EXISTS article_search_idx ON articles USING GIN (search_vector);
	`
	if _, err := dbPool.Exec(context.Background(), createTableSQL); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Repositories & Services
	newsRepo := repository.NewNewsRepository(dbPool)

	// Seed Data
	ingester := service.NewIngestionService(newsRepo)
	go func() {
		if err := ingester.IngestFromFile(context.Background(), "data/news_data.json"); err != nil {
			logger.Log.Error("Seeding failed (might be already seeded or file missing)", "error", err)
		}
	}()

	// HTTP Server
	r := gin.Default()

	logger.Log.Info("Server starting on port " + cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
