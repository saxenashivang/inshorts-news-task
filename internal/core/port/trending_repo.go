package port

import "context"

type TrendingRepository interface {
	RecordView(ctx context.Context, articleID int64) error
	GetTrending(ctx context.Context, limit int64) ([]int64, error)
}
