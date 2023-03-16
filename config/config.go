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

	/*	if err := cfg.Validate(); err != nil {
		return nil, err
	}*/

	return &cfg, nil
}

const (
	ModePro = "prod"
	ModeDev = "dev"
)

const appName = "xapi"

func DefaultConfig() Config {
	return Config{
		App:  appName,
		Mode: ModeDev,
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
	App        string      `mapstructure:"app"`
	Mode       string      `mapstructure:"mode"`
	Serve      *Servers    `mapstructure:"serve"`
	Cache      *Cache      `mapstructure:"cache"`
	Pubsub     *Pubsub     `mapstructure:"pubsub"`
	Logger     *Logger     `mapstructure:"logger"`
	Queue      *Queue      `mapstructure:"queue"`
	DB         *DB         `mapstructure:"db"`
	Prometheus *Prometheus `mapstructure:"prometheus"`
}

func (c Config) Validate() error {
	parts := []interface{}{
		c.Serve,
		c.Cache,
		c.Pubsub,
		c.Logger,
		c.Queue,
		c.DB,
		c.Prometheus,
	}

	for _, i := range parts {
		if part, ok := i.(interface{ Validate() error }); ok {
			if err := part.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
