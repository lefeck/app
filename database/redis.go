package database

import (
	"app/config"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	RedisDisableError = errors.New("redis disable")
)

type RedisDB struct {
	enable bool
	*redis.Client
}

func NewRedis(conf *config.RedisConfig) (*RedisDB, error) {
	if !conf.Enable {
		logrus.Info("redis disable")
		return &RedisDB{}, nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	})
	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}
	return &RedisDB{
		enable: true,
		Client: client,
	}, nil
}

func (rdb *RedisDB) Enable() bool {
	return rdb.enable
}

func (rdb *RedisDB) HGet(key, field string, obj interface{}) error {
	if !rdb.enable {
		return RedisDisableError
	}
	return rdb.Client.HGet(context.Background(), key, field).Scan(obj)
}

func (rdb *RedisDB) HSet(key, field string, val interface{}) error {
	if !rdb.enable {
		return RedisDisableError
	}
	return rdb.Client.HSet(context.Background(), key, field, val).Err()
}

func (rdb *RedisDB) HDel(key string, field ...string) error {
	if !rdb.enable {
		return RedisDisableError
	}
	return rdb.Client.HDel(context.Background(), key, field...).Err()
}
