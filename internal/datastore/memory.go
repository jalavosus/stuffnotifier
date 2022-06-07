package datastore

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/dgraph-io/ristretto"
	ristrettoz "github.com/dgraph-io/ristretto/z"
	"github.com/pkg/errors"

	"github.com/jalavosus/stuffnotifier/internal/utils"
)

const (
	inMemoryCacheMaxItems    = 3_000
	inMemoryCacheNumCounters = inMemoryCacheMaxItems * 10 // 30_000
	inMemoryCacheMaxCost     = 1_000
	inMemoryCacheBufferItems = 64
	inMemoryCacheItemCost    = 1
)

type memoryDatastore[T any] struct {
	*baseDatastore[T]
	client       *ristretto.Cache
	clientConfig *ristretto.Config
}

func NewInMemoryDatastore[T any](conf *Config) (Datastore[T], error) {
	var err error

	datastoreConf := defaultDatastoreConfig()
	if conf != nil {
		datastoreConf = conf.datastoreConfig()
	}

	m := &memoryDatastore[T]{
		baseDatastore: newBaseDatastore[T](datastoreConf),
		clientConfig: &ristretto.Config{
			NumCounters: inMemoryCacheNumCounters,
			MaxCost:     inMemoryCacheMaxCost,
			BufferItems: inMemoryCacheBufferItems,
			KeyToHash:   keyHasher,
		},
	}

	m.client, err = ristretto.NewCache(m.clientConfig)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *memoryDatastore[T]) Exists(_ context.Context, key string) (bool, error) {
	_, ok := m.client.Get(m.prefixKey(key))
	return ok, nil
}

func (m *memoryDatastore[T]) CheckTtl(_ context.Context, key string) (time.Duration, bool, error) {
	ttl, ok := m.client.GetTTL(m.prefixKey(key))
	return ttl, ok, nil
}

func (m *memoryDatastore[T]) Get(_ context.Context, key string) (*T, bool, error) {
	res, ok := m.client.Get(m.prefixKey(key))
	if !ok {
		return nil, false, nil
	}

	return fullDecode[T](res.(string))
}

func (m *memoryDatastore[T]) Insert(_ context.Context, key string, data T) error {
	encoded, err := fullEncode[T](data)
	if err != nil {
		return err
	}

	ok := m.client.SetWithTTL(m.prefixKey(key), encoded, inMemoryCacheItemCost, m.DefaultTtl())
	if !ok {
		return errors.Errorf("unable to add data with key %[1]s to cache", key)
	}

	return nil
}

func (m *memoryDatastore[T]) Delete(_ context.Context, key string) error {
	m.client.Del(m.prefixKey(key))
	return nil
}

func (m *memoryDatastore[T]) UpdateTtl(_ context.Context, key string, newTtl time.Duration) (bool, error) {
	key = m.prefixKey(key)

	res, ok := m.client.Get(key)
	if !ok {
		return false, nil
	}

	ok = m.client.SetWithTTL(key, res, inMemoryCacheItemCost, newTtl)
	if !ok {
		return false, errors.Errorf("unable to update ttl for key %[1]s", key)
	}

	return true, nil
}

func keyHasher(key any) (x, y uint64) {
	k, ok := key.(string)
	if !ok {
		return ristrettoz.KeyToHash(key)
	}

	hash := utils.SHA256(k)

	hashSplitLen := len(hash) / 2
	x = binary.LittleEndian.Uint64(hash[:hashSplitLen])
	y = binary.LittleEndian.Uint64(hash[hashSplitLen:])

	return
}
