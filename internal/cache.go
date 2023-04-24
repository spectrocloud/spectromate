// Copyright (c) Spectro Cloud
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Cache is an interface for a cache. The implementation is up to the user.
// The default implementation is Redis.
type Cache interface {
	StoreHashMap(ctx context.Context, primaryKey string, item map[string]interface{}) error
	GetHashMap(ctx context.Context, primaryKey string) (bool, map[string]string, error)
	ExpireKey(ctx context.Context, key string, t time.Duration) error
	Ping() error
}

func NewCache(connectionString, redisUser, redisPassword string, tls *tls.Config) Cache {
	return &RedisCache{
		redis: redis.NewClient(&redis.Options{
			Addr:      connectionString,
			Password:  redisPassword,
			DB:        0, // use default DB
			Username:  redisUser,
			TLSConfig: tls,
		}),
	}
}

// RedisCache is a Redis implementation of the Cache interface.
type RedisCache struct {
	redis *redis.Client
}

// StoreHashMap stores a hash map in the database.
func (c *RedisCache) StoreHashMap(ctx context.Context, primaryKey string, item map[string]interface{}) error {

	err := c.redis.HSet(ctx, primaryKey, item).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error storing item entry in cache.")
		return err
	}

	return nil
}

// ExpireKey sets an expiration on a cache key.
func (c *RedisCache) ExpireKey(ctx context.Context, key string, t time.Duration) error {

	err := c.redis.Expire(ctx, key, t).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error setting expiration on cache key")
		return err
	}

	return nil
}

// GetHashMap gets a hash map from the database.
func (c *RedisCache) GetHashMap(ctx context.Context, primaryKey string) (bool, map[string]string, error) {

	result, err := c.redis.HGetAll(ctx, primaryKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug().Msgf("key not found in cache: %s", primaryKey)
			return false, nil, nil
		}
		log.Error().Err(err).Msg("error retrieving key from cache.")
		return false, nil, err
	}

	if len(result) == 0 {
		log.Debug().Msgf("key is empty: %s", primaryKey)
		return false, nil, nil
	}

	return true, result, nil
}

// Ping checks the connection to the database.
func (r *RedisCache) Ping() error {
	return r.redis.Ping(context.Background()).Err()
}
