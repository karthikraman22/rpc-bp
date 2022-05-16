package config

import (
	"strings"

	"github.com/karthikraman22/rpc-bp/logger"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	*koanf.Koanf
}

func NewConfig(fname string, envPrefix string) (*Config, error) {
	log := logger.WithName("config")
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var conf = koanf.New(".")

	log.Info("loading configurations", "config_file", fname, "env_prefix", envPrefix)

	// Load JSON config.
	if err := conf.Load(file.Provider(fname), yaml.Parser()); err != nil {
		log.Error(err, "config_file_load_error")
		return nil, err
	}

	err := conf.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	if err != nil {
		log.Error(err, "config_env_load_error")
	}
	return &Config{conf}, err
}
