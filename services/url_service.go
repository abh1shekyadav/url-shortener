package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/abh1shekyadav/url-shortener/repositories"
	"github.com/abh1shekyadav/url-shortener/utils"
)

type URLService struct {
	repo repositories.URLRepository
}

func NewURLService(repo repositories.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string) (string, *repositories.URLData, error) {
	var code string
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		code = utils.GenerateShortCode()
		log.Printf("[SERVICE] Attempt %d: Generated shortcode=%s", i+1, code)

		exists, err := s.repo.IsCodeExists(ctx, code)
		if err != nil {
			log.Printf("[SERVICE] Error checking code existence (%s): %v", code, err)
			return "", nil, fmt.Errorf("check existence failed: %w", err)
		}

		if !exists {
			log.Printf("[SERVICE] Code %s is unique, proceeding to save", code)
			break
		}

		log.Printf("[SERVICE] Collision detected for %s, regenerating...", code)
		code = ""
	}

	if code == "" {
		return "", nil, fmt.Errorf("failed to generate unique shortcode after %d attempts", maxRetries)
	}

	data := &repositories.URLData{
		ShortCode:   code,
		OriginalURL: originalURL,
		ExpiresAt:   s.getExpiryTime(),
	}

	if err := s.repo.Save(ctx, data); err != nil {
		log.Printf("[SERVICE] Failed to save URL %s (%s): %v", originalURL, code, err)
		return "", nil, fmt.Errorf("save failed: %w", err)
	}

	log.Printf("[SERVICE] ✅ Saved URL=%s as code=%s expires_at=%v", originalURL, code, data.ExpiresAt)
	return code, data, nil
}

func (s *URLService) ResolveURL(ctx context.Context, code string) (string, error) {
	log.Printf("[SERVICE] Resolving shortcode=%s", code)
	url, err := s.repo.IncrementClick(ctx, code)
	if err != nil {
		log.Printf("[SERVICE] Error resolving shortcode %s: %v", code, err)
	}
	return url, err
}

func (s *URLService) GetStats(ctx context.Context, code string) (*repositories.URLData, error) {
	log.Printf("[SERVICE] Fetching stats for %s", code)
	data, err := s.repo.GetStats(ctx, code)
	if err != nil {
		log.Printf("[SERVICE] Error fetching stats for %s: %v", code, err)
	}
	return data, err
}

func (s *URLService) getExpiryTime() time.Time {
	days, err := strconv.Atoi(os.Getenv("URL_EXPIRY_DAYS"))
	if err != nil || days <= 0 {
		days = 7
	}
	return time.Now().Add(time.Hour * 24 * time.Duration(days))
}
