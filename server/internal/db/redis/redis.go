package redis

import (
	"context"
	"degrens/panel/internal/config"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

var Redis *Client

func InitRedisClient(conf *config.Config) *Client {
	Redis = &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
			Password: conf.Redis.Password,
			DB:       0,
		}),
		ctx: context.Background(),
	}
	err := Redis.client.Ping(Redis.ctx).Err()
	if err != nil {
		logrus.Fatalf("Failed to connect to redis: %s", err)
	}

	logrus.Info("Connected to redis")
	return Redis
}

// Generate a UUID that has no infaormation stored in redis
func (r *Client) GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	// Check if there is no data in redis with this id
	for r.client.Exists(r.ctx, id.String()).Val() != 0 {
		id, err = uuid.NewUUID()
		if err != nil {
			return "", err
		}
	}
	return id.String(), nil
}

func (r *Client) Set(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(r.ctx, key, p, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) Remove(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}
