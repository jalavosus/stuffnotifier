package datastore

import (
	"time"

	"github.com/jalavosus/stuffnotifier/internal/utils"
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

type RedisDatastoreConfig struct {
	Host       *string        `json:"host,omitempty" yaml:"host,omitempty" toml:"Host,omitempty"`
	Port       *int           `json:"port,omitempty" yaml:"port,omitempty" toml:"Port,omitempty"`
	Password   *string        `json:"password,omitempty" yaml:"password,omitempty" toml:"Password,omitempty"`
	Prefix     *string        `json:"prefix,omitempty" yaml:"prefix,omitempty" toml:"Prefix,omitempty"`
	DefaultTtl *time.Duration `json:"default_ttl,omitempty" yaml:"default_ttl,omitempty" toml:"DefaultTTL,omitempty"`
}

func (c *RedisDatastoreConfig) toInternalConfig() redisDatastoreConfig {
	config := defaultRedisDatastoreConfig()

	if confHost, ok := utils.FromPointer(c.Host); ok && confHost != "" {
		config.Host = confHost
	}

	if confPort, ok := utils.FromPointer(c.Port); ok && confPort != 0 {
		config.Port = confPort
	}

	if confPass, ok := utils.FromPointer(c.Password); ok {
		config.Password = confPass
	}

	if confPrefix, ok := utils.FromPointer(c.Prefix); ok {
		config.Prefix = confPrefix
	}

	if confTtl, ok := utils.FromPointer(c.DefaultTtl); ok && confTtl != time.Duration(0) {
		config.DefaultTtl = confTtl
	}

	return config
}

type redisDatastoreConfig struct {
	Host       string
	Password   string
	Prefix     string
	Port       int
	DefaultTtl time.Duration
}

func defaultRedisDatastoreConfig() redisDatastoreConfig {
	return redisDatastoreConfig{
		Host:       redisDefaultHost,
		Port:       redisDefaultPort,
		Password:   redisDefaultPassword,
		Prefix:     DefaultKeyPrefix,
		DefaultTtl: DefaultTtl,
	}
}
