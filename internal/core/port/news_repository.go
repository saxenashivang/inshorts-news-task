package port

import (
	"context"

	"github.com/shivangsaxena/inshorts-task/internal/core/entity"
)

type NewsRepository interface {
	BulkInsert(ctx context.Context, articles []entity.Article) error
	Search(ctx context.Context, query string, filters map[string]interface{}) ([]entity.Article, error)
	GetNearby(ctx context.Context, lat, lng, radius float64) ([]entity.Article, error)
}
