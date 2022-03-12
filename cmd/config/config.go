package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/time/rate"
)

type Config struct {
	RedisAddr           string        `envconfig:"REDIS_ADDR"`
	RedisMasterName     string        `envconfig:"REDIS_MASTER_NAME"`
	RedisPass           string        `envconfig:"REDIS_PASS"`
	RedisInstancesCount int           `envconfig:"REDIS_INSTANCES_COUNT"`
	RedisDialTimeout    time.Duration `envconfig:"REDIS_DIAL_TIMEOUT"`
	RedisReadTimeout    time.Duration `envconfig:"REDIS_READ_TIMEOUT"`
	RedisWriteTimeout   time.Duration `envconfig:"REDIS_WRITE_TIMEOUT"`
	RedisPoolSize       int           `envconfig:"REDIS_POOL_SIZE"`
	RedisMinIdleConn    int           `envconfig:"REDIS_MIN_IDLE_CONNECTIONS"`
	RedisExpiration     time.Duration `envconfig:"REDIS_EXPIRATION"`
	QuotesURL           string        `envconfig:"QUOTES_URL"`
	QuotesTimeout       time.Duration `envconfig:"QUOTES_TIMEOUT"`
	TargetBits          int           `envconfig:"TARGET_BITS"`
	RateLimit           rate.Limit    `envconfig:"RATE_LIMIT"`
	Burst               int           `envconfig:"BURST"`
}

func NewConfig() *Config {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, ignore outside the local")
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("envconfig err: %v", err.Error())
	}
	log.Println("envconfig ok")

	return cfg
}
