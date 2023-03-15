package config

import (
	"github.com/spf13/viper"
)

func New(v *viper.Viper) (*Config, error) {
	cfg := DefaultConfig()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() Config {
	return Config{
		Serve: &Servers{
			Public: &Server{
				Port: 11230,
			},
			Admin: &Server{
				Port: 11231,
			},
			Metrics: &Server{
				Port: 11232,
			},
			Healthz: &Server{
				Port: 11233,
			},
		},
	}
}

type Config struct {
	Serve  *Servers `mapstructure:"serve"`
	Cache  *Cache   `mapstructure:"cache"`
	Pubsub *Pubsub  `mapstructure:"pubsub"`
	Logger *Logger  `mapstructure:"logger"`
	Queue  *Queue   `mapstructure:"queue"`
}
