package repositories

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/abh1shekyadav/url-shortener/config"
)

type RedisURLRepository struct {
	primary URLRepository
}

func NewRedisURLRepository(primary URLRepository) *RedisURLRepository {
	return &RedisURLRepository{primary: primary}
}

func (repo *RedisURLRepository) Save(ctx context.Context, data *URLData) error {
	err := repo.primary.Save(ctx, data)
	if err != nil {
		return err
	}
	bytes, _ := json.Marshal(data)
	config.RedisClient.Set(ctx, data.ShortCode, bytes, 24*time.Hour)
	log.Printf("[CACHE] Saved %s to Redis\n", data.ShortCode)
	return nil
}

func (repo *RedisURLRepository) FindByCode(ctx context.Context, code string) (*URLData, error) {
	val, err := config.RedisClient.Get(ctx, code).Result()
	if err == nil {
		var data URLData
		if jsonErr := json.Unmarshal([]byte(val), &data); jsonErr == nil {
			log.Printf("[CACHE] Hit for %s\n", code)
			return &data, nil
		}
	}
	log.Printf("[CACHE] Miss for %s\n", code)
	data, dbErr := repo.primary.FindByCode(ctx, code)
	if dbErr == nil {
		bytes, _ := json.Marshal(data)
		config.RedisClient.Set(ctx, data.ShortCode, bytes, 24*time.Hour)
	}
	return data, dbErr
}

func (repo *RedisURLRepository) IncrementClick(ctx context.Context, code string) (string, error) {
	original, err := repo.primary.IncrementClick(ctx, code)
	if err == nil {
		config.RedisClient.Del(ctx, code)
	}
	return original, err
}

func (repo *RedisURLRepository) GetStats(ctx context.Context, code string) (*URLData, error) {
	return repo.primary.GetStats(ctx, code)
}

func (repo *RedisURLRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	return repo.primary.IsCodeExists(ctx, code)
}
