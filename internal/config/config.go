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
		Serve: Serve{
			Public: &ServeItem{
				Port: 11230,
			},
			Admin: &ServeItem{
				Port: 11231,
			},
			Metrics: &ServeItem{
				Port: 11232,
			},
			Healthz: &ServeItem{
				Port: 11233,
			},
		},
	}
}

type Config struct {
	Serve  *Serve  `mapstructure:"serve"`
	Cache  *Cache  `mapstructure:"cache"`
	Pubsub *Pubsub `mapstructure:"pubsub"`
	Logger *Logger `mapstructure:"logger"`
	Queue  *Queue  `mapstructure:"queue"`
}
