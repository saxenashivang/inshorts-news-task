package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type TrendingRepo struct {
	client *redis.Client
}

func NewTrendingRepo(addr string) *TrendingRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}
	return &TrendingRepo{client: rdb}
}

func (r *TrendingRepo) RecordView(ctx context.Context, articleID string) error {
	// Increment views
	key := fmt.Sprintf("article:%s:views", articleID)
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// Calculate Trending Score
	// Simple Logic: views + (now - published_approx / decay_factor)
	score := float64(val)

	return r.client.ZAdd(ctx, "trending:global", redis.Z{
		Score:  score,
		Member: articleID,
	}).Err()
}

func (r *TrendingRepo) GetTrending(ctx context.Context, limit int64) ([]string, error) {
	// Get top articles
	vals, err := r.client.ZRevRange(ctx, "trending:global", 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	return vals, nil
}
