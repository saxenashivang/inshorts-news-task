package usecase

import (
	"context"

	"github.com/shivangsaxena/inshorts-task/internal/core/entity"
	"github.com/shivangsaxena/inshorts-task/internal/core/port"
)

type NewsUseCase struct {
	repo port.NewsRepository
}

func NewNewsUseCase(repo port.NewsRepository) *NewsUseCase {
	return &NewsUseCase{repo: repo}
}

type NewsResponse struct {
	Articles []entity.Article `json:"articles"`
	Summary  string           `json:"ai_summary,omitempty"`
}

func (uc *NewsUseCase) GetNews(ctx context.Context, query string, userLat, userLng float64) (*NewsResponse, error) {
	return &NewsResponse{
		Articles: []entity.Article{},
		Summary:  "",
	}, nil
}
