package repositories

import (
	"context"
	"log"

	"github.com/abh1shekyadav/url-shortener/config"
)

type PostgresURLRepository struct{}

func NewPostgresURLRepository() *PostgresURLRepository {
	return &PostgresURLRepository{}
}

func (r *PostgresURLRepository) Save(ctx context.Context, data *URLData) error {
	log.Printf("[REPO] Inserting shortcode=%s URL=%s\n", data.ShortCode, data.OriginalURL)
	_, err := config.DB.Exec(ctx,
		"INSERT INTO urls (short_code, original_url) VALUES ($1, $2)",
		data.ShortCode, data.OriginalURL)
	if err != nil {
		log.Printf("[REPO] Insert failed for code=%s: %v\n", data.ShortCode, err)
	} else {
		log.Printf("[REPO] Insert successful for code=%s\n", data.ShortCode)
	}
	return err
}

func (r *PostgresURLRepository) FindByCode(ctx context.Context, code string) (*URLData, error) {
	log.Printf("[REPO] Looking up code=%s\n", code)
	row := config.DB.QueryRow(ctx,
		"SELECT id, short_code, original_url, click_count FROM urls WHERE short_code=$1", code)

	var u URLData
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount); err != nil {
		log.Printf("[REPO] Code=%s not found: %v\n", code, err)
		return nil, err
	}
	return &u, nil
}

func (r *PostgresURLRepository) IncrementClick(ctx context.Context, code string) (string, error) {
	log.Printf("[REPO] Incrementing click count for code=%s\n", code)
	var original string
	err := config.DB.QueryRow(ctx,
		"UPDATE urls SET click_count = click_count + 1 WHERE short_code=$1 RETURNING original_url", code).
		Scan(&original)
	if err != nil {
		log.Printf("[REPO] Error incrementing click count for %s: %v\n", code, err)
	}
	return original, err
}

func (r *PostgresURLRepository) GetStats(ctx context.Context, code string) (*URLData, error) {
	log.Printf("[REPO] Getting stats for code=%s\n", code)
	row := config.DB.QueryRow(ctx,
		"SELECT id, short_code, original_url, click_count FROM urls WHERE short_code=$1", code)

	var u URLData
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount); err != nil {
		log.Printf("[REPO] Error fetching stats for %s: %v\n", code, err)
		return nil, err
	}
	return &u, nil
}

func (r *PostgresURLRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code=$1)", code).Scan(&exists)
	log.Printf("[REPO] Checked existence for code=%s -> exists=%v err=%v\n", code, exists, err)
	return exists, err
}
