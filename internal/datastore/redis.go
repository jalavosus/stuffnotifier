package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

const (
	redisDefaultHost     string = "localhost"
	redisDefaultPort     int    = 6379
	redisDefaultPassword string = ""
)

type redisDatastore[T any] struct {
	*baseDatastore[T]
	client       *redis.Client
	clientConfig redisDatastoreConfig
}

// NewRedisDatastore wraps NewRedisDatastoreContext,
// passing a new context.Background() context.
func NewRedisDatastore[T any](conf *Config) (Datastore[T], error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return NewRedisDatastoreContext[T](ctx, conf)
}

// NewRedisDatastoreContext returns a Datastore implementation which
// uses Redis as the underlying datastore.
func NewRedisDatastoreContext[T any](ctx context.Context, conf *Config) (Datastore[T], error) {
	datastoreConf := defaultDatastoreConfig()
	clientConfig := defaultRedisDatastoreConfig()

	if conf != nil {
		datastoreConf = conf.datastoreConfig()

		if conf.Redis != nil {
			clientConfig = conf.Redis.toInternalConfig()
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%[1]s:%[2]d", clientConfig.Host, clientConfig.Port),
		Password: clientConfig.Password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "error pinging redis server")
	}

	return &redisDatastore[T]{
		baseDatastore: newBaseDatastore[T](datastoreConf),
		client:        client,
		clientConfig:  clientConfig,
	}, nil
}

func (d *redisDatastore[T]) Exists(ctx context.Context, key string) (bool, error) {
	res, err := d.client.Exists(ctx, d.prefixKey(key)).Result()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (d *redisDatastore[T]) CheckTtl(ctx context.Context, key string) (time.Duration, bool, error) {
	res, err := d.client.TTL(ctx, d.prefixKey(key)).Result()
	if err != nil {
		return time.Duration(0), false, err
	}

	if res == -1 {
		return time.Duration(0), true, nil
	} else if res == -2 {
		return time.Duration(0), false, nil
	}

	return res, true, nil
}

func (d *redisDatastore[T]) Get(ctx context.Context, key string) (*T, bool, error) {
	res, err := d.client.Get(ctx, d.prefixKey(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}

		return nil, false, err
	}

	return fullDecode[T](res)
}

func (d *redisDatastore[T]) Insert(ctx context.Context, key string, data T) error {
	encoded, err := fullEncode[T](data)
	if err != nil {
		return err
	}

	return d.client.Set(ctx, d.prefixKey(key), encoded, d.DefaultTtl()).Err()
}

func (d *redisDatastore[T]) Delete(ctx context.Context, key string) error {
	_, err := d.client.Del(ctx, d.prefixKey(key)).Result()
	return err
}

func (d *redisDatastore[T]) UpdateTtl(ctx context.Context, key string, newTtl time.Duration) (bool, error) {
	return d.client.Expire(ctx, d.prefixKey(key), newTtl).Result()
}
