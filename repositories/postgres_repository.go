package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/abh1shekyadav/url-shortener/config"
	"github.com/jackc/pgx/v5"
)

const (
	insertURL = `
		INSERT INTO urls (short_code, original_url, expires_at)
		VALUES ($1, $2, $3)`
	selectByCode = `
		SELECT id, short_code, original_url, click_count, expires_at
		FROM urls
		WHERE short_code=$1 AND expires_at > NOW()`
	updateClicks = `
		UPDATE urls
		SET click_count = click_count + 1
		WHERE short_code=$1 AND expires_at > NOW()
		RETURNING original_url`
	selectStats = `
		SELECT id, short_code, original_url, click_count, expires_at
		FROM urls
		WHERE short_code=$1`
	checkExists = `
		SELECT EXISTS(SELECT 1 FROM urls WHERE short_code=$1)`
)

type PostgresURLRepository struct{}

func NewPostgresURLRepository() *PostgresURLRepository {
	return &PostgresURLRepository{}
}

func (r *PostgresURLRepository) Save(ctx context.Context, data *URLData) error {
	log.Printf("[REPO] Inserting shortcode=%s URL=%s expiry=%s", data.ShortCode, data.OriginalURL, data.ExpiresAt)
	_, err := config.DB.Exec(ctx, insertURL, data.ShortCode, data.OriginalURL, data.ExpiresAt)
	if err != nil {
		log.Printf("[REPO] Insert failed for code=%s: %v", data.ShortCode, err)
	}
	return err
}

func (r *PostgresURLRepository) FindByCode(ctx context.Context, code string) (*URLData, error) {
	log.Printf("[REPO] Looking up code=%s", code)
	row := config.DB.QueryRow(ctx, selectByCode, code)
	var u URLData
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[REPO] Code=%s not found", code)
			return nil, nil
		}
		log.Printf("[REPO] Query failed for code=%s: %v", code, err)
		return nil, err
	}
	return &u, nil
}

func (r *PostgresURLRepository) IncrementClick(ctx context.Context, code string) (string, error) {
	log.Printf("[REPO] Incrementing click count for code=%s", code)
	var original string
	err := config.DB.QueryRow(ctx, updateClicks, code).Scan(&original)
	if err != nil {
		log.Printf("[REPO] Error incrementing click count for %s: %v", code, err)
	}
	return original, err
}

func (r *PostgresURLRepository) GetStats(ctx context.Context, code string) (*URLData, error) {
	log.Printf("[REPO] Getting stats for code=%s", code)
	row := config.DB.QueryRow(ctx, selectStats, code)
	var u URLData
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.ExpiresAt); err != nil {
		log.Printf("[REPO] Error fetching stats for %s: %v", code, err)
		return nil, err
	}
	return &u, nil
}

func (r *PostgresURLRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow(ctx, checkExists, code).Scan(&exists)
	log.Printf("[REPO] Checked existence for code=%s -> exists=%v err=%v", code, exists, err)
	return exists, err
}
