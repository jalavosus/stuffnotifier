package datastore

import (
	"time"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
)

const (
	// DefaultTtl is the default time-to-live for all datastore implementations
	DefaultTtl = (8 * time.Hour) + (30 * time.Minute)
	// DefaultKeyPrefix is the key prefix used by datastore implementations
	// when one is not manually configured.
	DefaultKeyPrefix = "lavamonster"
)

type Config struct {
	KeyPrefix  *string               `json:"key_prefix,omitempty" yaml:"key_prefix,omitempty" toml:"KeyPrefix,omitempty"`
	DefaultTtl *time.Duration        `json:"default_ttl,omitempty" yaml:"default_ttl,omitempty" toml:"DefaultTtl,omitempty"`
	Redis      *RedisDatastoreConfig `json:"redis,omitempty" yaml:"redis,omitempty" toml:"Redis,omitempty"`
}

func (c Config) datastoreConfig() datastoreConfig {
	conf := datastoreConfig{
		KeyPrefix:  DefaultKeyPrefix,
		DefaultTtl: DefaultTtl,
	}

	if keyPrefix, ok := utils.FromPointer(c.KeyPrefix); ok {
		conf.KeyPrefix = keyPrefix
	}

	if defaultTtl, ok := utils.FromPointer(c.DefaultTtl); ok && defaultTtl != time.Duration(0) {
		conf.DefaultTtl = defaultTtl
	}

	return conf
}

type datastoreConfig struct {
	KeyPrefix  string
	DefaultTtl time.Duration
}

func defaultDatastoreConfig() datastoreConfig {
	return datastoreConfig{
		KeyPrefix:  DefaultKeyPrefix,
		DefaultTtl: DefaultTtl,
	}
}
