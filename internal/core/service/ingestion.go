package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/shivangsaxena/inshorts-task/internal/core/entity"
	"github.com/shivangsaxena/inshorts-task/internal/core/port"
)

type IngestionService struct {
	repo port.NewsRepository
}

func NewIngestionService(repo port.NewsRepository) *IngestionService {
	return &IngestionService{repo: repo}
}

func (s *IngestionService) IngestFromFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	// Define a DTO to match JSON structure exactly
	type articleDTO struct {
		ID        string   `json:"id"`
		Title     string   `json:"title"`
		Content   string   `json:"description"`
		Source    string   `json:"source_name"`
		URL       string   `json:"url"`
		Category  []string `json:"category"`
		Lat       float64  `json:"latitude"`
		Lng       float64  `json:"longitude"`
		Published string   `json:"publication_date"`
		Score     float64  `json:"relevance_score"`
	}

	var dtos []articleDTO
	if err := json.Unmarshal(byteValue, &dtos); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	// Convert DTOs to Entities
	var articles []entity.Article
	const layout = "2006-01-02T15:04:05"

	for _, dto := range dtos {
		pubTime, err := time.Parse(layout, dto.Published)
		if err != nil {
			if dto.Published != "" {
				fmt.Printf("Warning: failed to parse time '%s' for article %s: %v\n", dto.Published, dto.ID, err)
			}
			pubTime = time.Now()
		}

		articles = append(articles, entity.Article{
			ID:        dto.ID,
			Title:     dto.Title,
			Content:   dto.Content,
			Source:    dto.Source,
			URL:       dto.URL,
			Category:  dto.Category,
			Lat:       dto.Lat,
			Lng:       dto.Lng,
			Published: pubTime,
			Score:     dto.Score,
		})
	}

	if err := s.repo.BulkInsert(ctx, articles); err != nil {
		return fmt.Errorf("failed to bulk insert: %w", err)
	}

	return nil
}
