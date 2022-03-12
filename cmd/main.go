package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"

	"pow-b/cmd/config"
	"pow-b/internal/app/requestchallenge"
	"pow-b/internal/app/solve"
	"pow-b/internal/middlewares"
	"pow-b/internal/pkg/auth"
	"pow-b/internal/pkg/quotes"
	solveService "pow-b/internal/pkg/solve"
	"pow-b/internal/pkg/storage/redis"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	spew.Dump(cfg)
	redisService, err := initRedis(ctx, cfg)
	if err != nil {
		log.Fatalf("couldn't init redis: %w", err)
	}

	authService := auth.New(redisService)
	reqChService := requestchallenge.New(authService)

	quotesService, err := initQuotes(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	solveImpl := initSolve(quotesService, cfg)

	public := chi.NewRouter()

	public.Use(middlewares.WithRateLimiter(cfg.RateLimit, cfg.Burst))

	public.Get("/request-challenge", func(writer http.ResponseWriter, request *http.Request) {
		reqChService.Gen(writer, request)
	})
	public.With(middlewares.WithAuth(authService)).Post("/solve-challenge", func(writer http.ResponseWriter, request *http.Request) {
		solveImpl.Solve(writer, request)
	})

	if err = http.ListenAndServe(":8081", public); err != nil {
		log.Fatal(ctx, err)
	}
}

func initRedis(ctx context.Context, cfg *config.Config) (*redis.Redis, error) {
	return redis.New(ctx, cfg.RedisAddr, cfg.RedisDialTimeout, cfg.RedisReadTimeout, cfg.RedisWriteTimeout, cfg.RedisPoolSize, cfg.RedisMinIdleConn, cfg.RedisExpiration)
}

func initQuotes(ctx context.Context, cfg *config.Config) (*quotes.Service, error) {
	quotesService := quotes.New(cfg.QuotesURL, cfg.QuotesTimeout)

	_, err := quotesService.Get(ctx)
	if err != nil {
		return &quotes.Service{}, fmt.Errorf("couldn't init quotes service: %w", err)
	}

	return quotesService, nil
}

func initSolve(q *quotes.Service, cfg *config.Config) *solve.Implemetation {
	service := solveService.New(cfg.TargetBits, q)

	return solve.New(service)
}
