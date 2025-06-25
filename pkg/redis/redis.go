package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr     string `json:"addr" yaml:"addr" toml:"addr" mapstructure:"addr"`
	Password string `json:"password" yaml:"password" toml:"password" mapstructure:"password"`
	DB       int    `json:"db" yaml:"db" toml:"db" mapstructure:"db"`
}

type Redis struct {
	*redis.Client
}

func New(config *Config) *Redis {

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // "" means no password set
		DB:       config.DB,       // 0 for using default DB
	})
	return &Redis{
		rdb,
	}
}

func (r *Redis) Pong(ctx context.Context) error {
	_, err := r.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
