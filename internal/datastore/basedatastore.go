package datastore

import (
	"time"
)

type baseDatastore[T any] struct {
	config datastoreConfig
}

func newBaseDatastore[T any](config datastoreConfig) *baseDatastore[T] {
	return &baseDatastore[T]{config}
}

func (d *baseDatastore[T]) DefaultTtl() time.Duration {
	return d.config.DefaultTtl
}

func (d *baseDatastore[T]) SetDefaultTtl(newTtl time.Duration) {
	d.config.DefaultTtl = newTtl
}

func (d *baseDatastore[T]) KeyPrefix() string {
	return d.config.KeyPrefix
}

func (d *baseDatastore[T]) SetKeyPrefix(newPrefix string) {
	d.config.KeyPrefix = newPrefix
}

func (d *baseDatastore[T]) prefixKey(key string) string {
	return d.KeyPrefix() + ":" + key
}
