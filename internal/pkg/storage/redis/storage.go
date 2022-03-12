package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Redis is a base struct
type Redis struct {
	client     *redis.Client
	expiration time.Duration
}

// New creates new client and returns error
func New(ctx context.Context, address string, dialTimeout, readTimeout, writeTimeout time.Duration, poolSize, minIdleConns int, expiration time.Duration) (*Redis, error) {
	r := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     "",
		DB:           0,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	if _, err := r.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("couldn't ping redis: %w", err)
	}

	client := new(Redis)
	client.client = r
	client.expiration = expiration

	return client, nil
}

// Add ...
func (r *Redis) Add(ctx context.Context, ip string, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis/Add")
	defer span.Finish()

	if _, err := r.client.Set(ctx, ip, token, r.expiration).Result(); err != nil {
		return err
	}

	return nil
}

// Check ...
func (r *Redis) Check(ctx context.Context, ip string, token string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis/Add")
	defer span.Finish()

	data, err := r.client.Get(ctx, ip).Result()
	if err != nil {
		log.Error(err)
		return false
	}

	return data == token
}

// Exist ...
func (r *Redis) Exist(ctx context.Context, ip string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis/Add")
	defer span.Finish()

	_, err := r.client.Get(ctx, ip).Result()

	return err == nil
}
