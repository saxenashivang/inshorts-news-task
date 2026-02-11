package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivangsaxena/inshorts-task/internal/core/entity"
)

type newsRepository struct {
	db *pgxpool.Pool
}

func NewNewsRepository(db *pgxpool.Pool) *newsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) BulkInsert(ctx context.Context, articles []entity.Article) error {
	batch := &pgx.Batch{}

	for _, a := range articles {
		// We insert geom using the PostGIS function directly
		sql := `INSERT INTO articles (
			id, title, content, source, url, category, lat, lng, published_at, relevance_score, geom
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, ST_SetSRID(ST_MakePoint($8, $7), 4326)::geography
		) ON CONFLICT (id) DO NOTHING`

		batch.Queue(sql,
			a.ID, a.Title, a.Content, a.Source, a.URL, a.Category, a.Lat, a.Lng, a.Published, a.Score,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	// Check for errors in the batch execution
	for i := 0; i < len(articles); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert article at index %d: %w", i, err)
		}
	}

	return nil
}

func (r *newsRepository) Search(ctx context.Context, query string, filters map[string]interface{}) ([]entity.Article, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	argIdx := 1

	sqlBuilder.WriteString("SELECT id, title, content, source, url, category, lat, lng, published_at, relevance_score FROM articles WHERE 1=1")

	if query != "" {
		fmt.Fprintf(&sqlBuilder, " AND search_vector @@ plainto_tsquery('english', $%d)", argIdx)
		args = append(args, query)
		argIdx++
	}

	if val, ok := filters["category"]; ok && val != "" {
		fmt.Fprintf(&sqlBuilder, " AND $%d = ANY(category)", argIdx)
		args = append(args, val)
		argIdx++
	}

	if val, ok := filters["source"]; ok && val != "" {
		fmt.Fprintf(&sqlBuilder, " AND source = $%d", argIdx)
		args = append(args, val)
		argIdx++
	}

	// Order by relevance and recency
	sqlBuilder.WriteString(" ORDER BY relevance_score DESC, published_at DESC LIMIT 50")

	rows, err := r.db.Query(ctx, sqlBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entity.Article])
}

func (r *newsRepository) GetNearby(ctx context.Context, lat, lng, radius float64) ([]entity.Article, error) {
	// radius in meters
	sql := `SELECT id, title, content, source, url, category, lat, lng, published_at, relevance_score 
			FROM articles 
			WHERE ST_DWithin(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
			ORDER BY ST_Distance(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography)
			LIMIT 50`

	rows, err := r.db.Query(ctx, sql, lng, lat, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entity.Article])
}
