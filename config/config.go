package config

import (
	"strings"

	"achuala.in/rpc-bp/logger"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	*koanf.Koanf
}

func NewConfig(fname string, envPrefix string) *Config {
	log := logger.WithName("config")
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var conf = koanf.New(".")

	// Load JSON config.
	if err := conf.Load(file.Provider(fname), yaml.Parser()); err != nil {
		log.Error(err, "error loading config")
	}

	conf.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	return &Config{conf}
}
