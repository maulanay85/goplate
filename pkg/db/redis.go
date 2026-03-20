package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	host     string
	port     int
	password string
	db       int

	// cluster
	clusterAddrs []string
}

type RedisOption func(*RedisConfig)

func WithRedisHost(host string) RedisOption {
	return func(c *RedisConfig) { c.host = host }
}

func WithRedisPort(port int) RedisOption {
	return func(c *RedisConfig) { c.port = port }
}

func WithRedisPassword(password string) RedisOption {
	return func(c *RedisConfig) { c.password = password }
}

func WithRedisDB(db int) RedisOption {
	return func(c *RedisConfig) { c.db = db }
}

// WithRedisClusterAddrs switches to cluster mode
func WithRedisClusterAddrs(addrs []string) RedisOption {
	return func(c *RedisConfig) { c.clusterAddrs = addrs }
}

// RedisClient is the interface the caller depends on
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration int) error
	Del(ctx context.Context, keys ...string) error
	Ping(ctx context.Context) error
	Close() error
}

// --- single node ---

type redisClient struct {
	client *redis.Client
}

func newSingleRedis(cfg *RedisConfig) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.host, cfg.port),
		Password: cfg.password,
		DB:       cfg.db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect single node: %w", err)
	}

	return &redisClient{client: client}, nil
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis: get %s: %w", key, err)
	}
	return val, nil
}

func (r *redisClient) Set(ctx context.Context, key string, value any, expiration int) error {
	if err := r.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("redis: set %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) Del(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis: del: %w", err)
	}
	return nil
}

func (r *redisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisClient) Close() error {
	return r.client.Close()
}

// --- cluster ---

type redisClusterClient struct {
	client *redis.ClusterClient
}

func newClusterRedis(cfg *RedisConfig) (RedisClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.clusterAddrs,
		Password: cfg.password,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect cluster: %w", err)
	}

	return &redisClusterClient{client: client}, nil
}

func (r *redisClusterClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis: get %s: %w", key, err)
	}
	return val, nil
}

func (r *redisClusterClient) Set(ctx context.Context, key string, value any, expiration int) error {
	if err := r.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("redis: set %s: %w", key, err)
	}
	return nil
}

func (r *redisClusterClient) Del(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis: del: %w", err)
	}
	return nil
}

func (r *redisClusterClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisClusterClient) Close() error {
	return r.client.Close()
}

func NewRedis(opts ...RedisOption) (RedisClient, error) {
	cfg := &RedisConfig{
		db: 0, // default
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// if cluster addrs are provided, use cluster mode
	if len(cfg.clusterAddrs) > 0 {
		return newClusterRedis(cfg)
	}

	return newSingleRedis(cfg)
}
