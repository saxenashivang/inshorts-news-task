package entity

import (
	"time"
)

type Article struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"description" db:"content"` // Maps to 'description' in JSON
	Source    string    `json:"source_name" db:"source"`  // Maps to 'source_name' in JSON
	URL       string    `json:"url" db:"url"`
	Category  []string  `json:"category" db:"category"`             // Maps to 'category' array in JSON
	Lat       float64   `json:"latitude" db:"lat"`                  // Maps to 'latitude'
	Lng       float64   `json:"longitude" db:"lng"`                 // Maps to 'longitude'
	Published time.Time `json:"publication_date" db:"published_at"` // Maps to 'publication_date'
	Score     float64   `json:"relevance_score" db:"relevance_score"`
	// PostGIS specific column, handled via custom expressions
}

func (Article) TableName() string {
	return "articles"
}
