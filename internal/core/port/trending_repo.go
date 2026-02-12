package port

import "context"

type TrendingRepository interface {
	RecordView(ctx context.Context, articleID string) error
	GetTrending(ctx context.Context, limit int64) ([]string, error)
}
