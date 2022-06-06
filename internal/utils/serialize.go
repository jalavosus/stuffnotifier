package utils

import (
	"encoding/json"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type ConfigFileType uint8

const (
	unknown ConfigFileType = iota
	JsonConfig
	YamlConfig
	TomlConfig
)

func UnmarshalConfig(data []byte, confType ConfigFileType, dest any) error {
	var unmarshalFn func([]byte, any) error

	switch confType {
	case JsonConfig:
		unmarshalFn = json.Unmarshal
	case YamlConfig:
		unmarshalFn = yaml.Unmarshal
	case TomlConfig:
		unmarshalFn = toml.Unmarshal
	default:
		return errors.Errorf("unknown file type")
	}

	return unmarshalFn(data, dest)
}

func ConfigFileTypeFromExtension(ext string) ConfigFileType {
	switch strings.ToLower(ext) {
	case "json", ".json":
		return JsonConfig
	case "yaml", ".yaml":
		return YamlConfig
	case "toml", ".toml":
		return TomlConfig
	default:
		return unknown
	}
}
