package middleware

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"time"
)

type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(key, value string, expiration time.Duration)
	Del(key string)
}

var CacheContextKey = "cache"

type Cache struct {
	l      hclog.Logger
	client *redis.Client
}

// Get retrieves the value for the given key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	// No context, give 10-seconds deadline.
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
		defer cancel()
	}

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.l.Info("Cache MISS", "key", key)
			return "", err
		}

		c.l.Error("Cache.Get failed with error", "error", err)
		return "", err
	}

	c.l.Info("Cache HIT", "url", key)
	return val, nil
}

// Set stores the given key:value for the given expiration duration.
func (c *Cache) Set(key, value string, expiration time.Duration) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	_, err := c.client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		c.l.Error("Cache.Set failed with error", "error", err)
		return
	}

	c.l.Info("Cache SET", "url", key)
}

func (c *Cache) Del(key string) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	c.client.Del(ctx, key)
}

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CacheContextKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseCaching(l hclog.Logger) *Cache {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		l.Error("Failed to connect to cache", "error", err)
		panic(err)
	}

	return &Cache{
		l:      l,
		client: client,
	}
}
