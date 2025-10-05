package services

import (
	"context"
	"log"

	"github.com/abh1shekyadav/url-shortener/repositories"
	"github.com/abh1shekyadav/url-shortener/utils"
)

type URLService struct {
	repo repositories.URLRepository
}

func NewURLService(repo repositories.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	var code string
	for {
		code = utils.GenerateShortCode()
		log.Println("[SERVICE] Generated shortcode:", code)

		exists, err := s.repo.IsCodeExists(ctx, code)
		if err != nil {
			log.Printf("[SERVICE] Error checking code existence (%s): %v\n", code, err)
			return "", err
		}
		if !exists {
			log.Println("[SERVICE] Code is unique, proceeding to save.")
			break
		}
		log.Println("[SERVICE] Collision detected, regenerating...")
	}

	err := s.repo.Save(ctx, &repositories.URLData{
		ShortCode:   code,
		OriginalURL: originalURL,
	})
	if err != nil {
		log.Printf("[SERVICE] Failed to save URL %s (%s): %v\n", originalURL, code, err)
		return "", err
	}

	log.Printf("[SERVICE] Successfully saved URL: %s as code: %s\n", originalURL, code)
	return code, nil
}

func (s *URLService) ResolveURL(ctx context.Context, code string) (string, error) {
	log.Println("[SERVICE] Resolving shortcode:", code)
	url, err := s.repo.IncrementClick(ctx, code)
	if err != nil {
		log.Printf("[SERVICE] Error resolving shortcode %s: %v\n", code, err)
	}
	return url, err
}

func (s *URLService) GetStats(ctx context.Context, code string) (*repositories.URLData, error) {
	log.Println("[SERVICE] Fetching stats for:", code)
	data, err := s.repo.GetStats(ctx, code)
	if err != nil {
		log.Printf("[SERVICE] Error fetching stats for %s: %v\n", code, err)
	}
	return data, err
}
