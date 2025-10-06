package repositories

import "context"

type URLData struct {
	ID          int
	ClickCount  int
	ShortCode   string
	OriginalURL string
}

type URLRepository interface {
	Save(ctx context.Context, data *URLData) error
	FindByCode(ctx context.Context, code string) (*URLData, error)
	IncrementClick(ctx context.Context, code string) (string, error)
	GetStats(ctx context.Context, code string) (*URLData, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
}
