package datastore

import (
	"context"
	"time"
)

type Datastore[T any] interface {
	// Exists checks if a key currently exists in the datastore.
	// Should return true if the key exists, false otherwise.
	// If the underlying datastore implementation returns an error
	// when checking key existence, Exists should return `false` and the
	// aforementioned error.
	Exists(ctx context.Context, key string) (bool, error)
	// CheckTtl returns the remaining time-to-live of data in the datastore
	// for the passed key.
	// If no data is set to the passed key, the zero-value of time.Duration is
	// returned, as well as false.
	// If the underlying datastore implementation returns an error,
	// that error will be returned.
	CheckTtl(ctx context.Context, key string) (time.Duration, bool, error)
	// Get returns the data which is stored with the passed key,
	// as well as a boolean representing whether the key exists
	// in the datastore at all.
	Get(ctx context.Context, key string) (*T, bool, error)
	// Insert inserts data using a given key into the datastore,
	// returning any error returned by the underlying datastore implementation.
	Insert(ctx context.Context, key string, data T) error
	// Delete removes a given key from the datastore, returning any error
	// returned by the underlying datastore implementation.
	// Note that this is a "blind" delete: if the passed key does not exist,
	// no error is returned, and no data is otherwise modified.
	Delete(ctx context.Context, key string) error
	// UpdateTtl sets the time-to-live for an object in the datastore which
	// is stored at the passed key.
	// If the passed key is unset, false is returned; otherwise, true is returned.
	UpdateTtl(ctx context.Context, key string, newTtl time.Duration) (bool, error)
	// DefaultTtl returns the time-to-live value used by default by the datastore
	// when storing new data.
	DefaultTtl() time.Duration
	// SetDefaultTtl sets the default time-to-live used by the datastore
	// when storing new data.
	SetDefaultTtl(newTtl time.Duration)
	// KeyPrefix returns the configured key prefixed used by the datastore.
	// Key prefixing should be used by all Datastore implementations
	// to ensure data consistency.
	KeyPrefix() string
	// SetKeyPrefix sets the key prefix used by the datastore
	// when storing new data.
	// Note that setting a new key prefix will *NOT* update previously stored
	// keys, essentially making data using the old key prefix inaccessible
	// to the datastore.
	SetKeyPrefix(newPrefix string)
}

func NewDatastore[T any](conf *Config) (Datastore[T], error) {
	if conf != nil && conf.Redis != nil {
		return NewRedisDatastore[T](conf)
	}

	return NewInMemoryDatastore[T](conf)
}
