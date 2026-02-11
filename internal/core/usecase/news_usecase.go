package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shivangsaxena/inshorts-task/internal/core/entity"
	"github.com/shivangsaxena/inshorts-task/internal/core/port"
)

type NewsUseCase struct {
	repo       port.NewsRepository
	llmService port.LLMService
}

func NewNewsUseCase(repo port.NewsRepository, llmService port.LLMService) *NewsUseCase {
	return &NewsUseCase{repo: repo, llmService: llmService}
}

type NewsResponse struct {
	Articles []entity.Article `json:"articles"`
	Summary  string           `json:"ai_summary,omitempty"`
}

func (uc *NewsUseCase) GetNews(ctx context.Context, query string, userLat, userLng float64) (*NewsResponse, error) {
	// 1. Analyze Intent
	parsed, err := uc.llmService.ParseQuery(ctx, query)
	if err != nil {
		// Fallback to text search if LLM fails
		parsed = &port.LLMResponse{Intent: "search"}
	}

	log.Printf("LLM Parsed Intent: %s, Location: %s, Category: %s, Entities: %v", parsed.Intent, parsed.Location, parsed.Category, parsed.Entities)

	var articles []entity.Article

	switch parsed.Intent {
	case "nearby":
		// use user location if parsed location is empty or "me"
		lat, lng := userLat, userLng
		// Mock geocoding for specific cities if needed (omitted for brevity)
		if strings.Contains(strings.ToLower(parsed.Location), "palo alto") {
			lat, lng = 37.4419, -122.1430
		}

		articles, err = uc.repo.GetNearby(ctx, lat, lng, 10000) // 10km radius
	default:
		// Search / Category
		filters := make(map[string]interface{})
		if parsed.Category != "" {
			filters["category"] = parsed.Category
		}
		// If explicit entities/keywords found, use them, else use original query
		searchQuery := query
		if len(parsed.Entities) > 0 {
			searchQuery = strings.Join(parsed.Entities, " ")
		}
		articles, err = uc.repo.Search(ctx, searchQuery, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}

	// 2. Enrich with Summary
	summary := ""
	if len(articles) > 0 {
		// Summarize the top 3 headlines to give context
		topHeadlines := []string{}
		for i := 0; i < 3 && i < len(articles); i++ {
			topHeadlines = append(topHeadlines, articles[i].Title)
		}
		summaryCtx := fmt.Sprintf("Query: %s. Top Headlines: %s", query, strings.Join(topHeadlines, "; "))
		summary, _ = uc.llmService.Summarize(ctx, summaryCtx)
	}

	return &NewsResponse{
		Articles: articles,
		Summary:  summary,
	}, nil
}
