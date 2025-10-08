package repositories

import (
	"context"
	"time"
)

type URLData struct {
	ID          int       `json:"id"`
	ClickCount  int       `json:"click_count"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type URLRepository interface {
	Save(ctx context.Context, data *URLData) error
	FindByCode(ctx context.Context, code string) (*URLData, error)
	IncrementClick(ctx context.Context, code string) (string, error)
	GetStats(ctx context.Context, code string) (*URLData, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
}
