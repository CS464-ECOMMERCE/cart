package services

import (
	"cart/configs"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps the Redis client with additional methods for cart operations
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	config configs.EnvConfig
}

var redisClient *RedisClient

// GetRedisClient returns a singleton instance of the Redis client
func GetRedisClient() *RedisClient {
	if redisClient == nil {
		config := configs.GetEnvConfig()
		ctx := context.Background()

		client := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		})

		// Test connection
		if _, err := client.Ping(ctx).Result(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}

		fmt.Println("Successfully connected to Redis")

		redisClient = &RedisClient{
			client: client,
			ctx:    ctx,
			config: config,
		}
	}

	return redisClient
}

// Close closes the Redis client connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// getCartKey returns the Redis key for a user's cart
func (r *RedisClient) getCartKey(userID uint64) string {
	return fmt.Sprintf("cart:%d", userID)
}

// SetCartTTL sets the TTL for a user's cart
func (r *RedisClient) SetCartTTL(userID uint64, ttl time.Duration) error {
	key := r.getCartKey(userID)
	return r.client.Expire(r.ctx, key, ttl).Err()
}

// GetCartTTL returns the TTL for a user's cart
func (r *RedisClient) GetCartTTL(userID uint64) (time.Duration, error) {
	key := r.getCartKey(userID)
	return r.client.TTL(r.ctx, key).Result()
}

// ExtendCartTTL extends the TTL for a user's cart
func (r *RedisClient) ExtendCartTTL(userID uint64) error {
	return r.SetCartTTL(userID, r.config.RedisDefaultTTL)
}

// DeleteCart deletes a user's cart
func (r *RedisClient) DeleteCart(userID uint64) error {
	key := r.getCartKey(userID)
	return r.client.Del(r.ctx, key).Err()
}

// ExecuteWithLock executes a function with a distributed lock
func (r *RedisClient) ExecuteWithLock(key string, ttl time.Duration, fn func() error) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	// Set lock with NX option (only set if key doesn't exist)
	ok, err := r.client.SetNX(r.ctx, lockKey, 1, ttl).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !ok {
		return fmt.Errorf("failed to acquire lock, already locked")
	}

	defer r.client.Del(r.ctx, lockKey)

	return fn()
}
