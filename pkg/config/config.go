package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Option func(*viper.Viper)

func Load[T any](opts ...Option) (T, error) {
	var cfg T

	v := viper.New()
	v.AutomaticEnv()

	for _, opt := range opts {
		opt(v)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("config: unmarshal failed: %w", err)
	}

	return cfg, nil
}

func WithPrefix(prefix string) Option {
	return func(v *viper.Viper) {
		v.SetEnvPrefix(prefix)
	}
}

func WithEnvFile(path string) Option {
	return func(v *viper.Viper) {
		v.SetConfigFile(path)
		v.SetConfigType("env")

		_ = v.ReadInConfig()
	}
}
