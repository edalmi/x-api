package config

import (
	"github.com/spf13/viper"
)

const (
	portPublic = iota + 11230
	portAdmin
	portMetricts
	portHealthz
)

func New(v *viper.Viper) (*Config, error) {
	cfg := DefaultConfig()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() Config {
	return Config{
		Serve: &Servers{
			Public: &Server{
				Port: portPublic,
			},
			Admin: &Server{
				Port: portAdmin,
			},
			Metrics: &Server{
				Port: portMetricts,
			},
			Healthz: &Server{
				Port: portHealthz,
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
	DB     DB       `mapstructure:"db"`
}

func (c Config) Validate() error {
	return nil
}
