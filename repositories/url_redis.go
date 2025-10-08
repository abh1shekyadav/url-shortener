package repositories

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/abh1shekyadav/url-shortener/config"
	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL     = 24 * time.Hour
	redisKeyPref = "url:"
)

type RedisURLRepository struct {
	primary URLRepository
}

func NewRedisURLRepository(primary URLRepository) *RedisURLRepository {
	return &RedisURLRepository{primary: primary}
}

func (repo *RedisURLRepository) Save(ctx context.Context, data *URLData) error {
	if err := repo.primary.Save(ctx, data); err != nil {
		return err
	}
	repo.cacheSet(ctx, data.ShortCode, data)
	return nil
}

func (repo *RedisURLRepository) FindByCode(ctx context.Context, code string) (*URLData, error) {
	val, err := config.RedisClient.Get(ctx, redisKeyPref+code).Result()
	if err != nil && err != redis.Nil {
		log.Printf("[CACHE] Redis unavailable, fallback to DB: %v", err)
		return repo.primary.FindByCode(ctx, code)
	}
	if err == nil {
		var data URLData
		if jsonErr := json.Unmarshal([]byte(val), &data); jsonErr != nil {
			log.Printf("[CACHE] Invalid cache entry for %s: %v", code, jsonErr)
		} else {
			log.Printf("[CACHE] Hit for %s\n", code)
			return &data, nil
		}
	}
	log.Printf("[CACHE] Miss for %s\n", code)
	data, dbErr := repo.primary.FindByCode(ctx, code)
	if dbErr == nil {
		repo.cacheSet(ctx, code, data)
	}
	return data, dbErr
}

func (repo *RedisURLRepository) IncrementClick(ctx context.Context, code string) (string, error) {
	original, err := repo.primary.IncrementClick(ctx, code)
	if err == nil {
		go func() {
			if delErr := config.RedisClient.Del(ctx, redisKeyPref+code).Err(); delErr != nil {
				log.Printf("[CACHE] Failed to delete key %s: %v", code, delErr)
			}
		}()
	}
	return original, err
}

func (repo *RedisURLRepository) GetStats(ctx context.Context, code string) (*URLData, error) {
	return repo.primary.GetStats(ctx, code)
}

func (repo *RedisURLRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	return repo.primary.IsCodeExists(ctx, code)
}

func (repo *RedisURLRepository) cacheSet(ctx context.Context, code string, data *URLData) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("[CACHE] Failed to marshal %s: %v", code, err)
		return
	}
	if err := config.RedisClient.Set(ctx, redisKeyPref+code, bytes, cacheTTL).Err(); err != nil {
		log.Printf("[CACHE] Failed to set key %s in Redis: %v", code, err)
	} else {
		log.Printf("[CACHE] Saved %s to Redis", code)
	}
}
