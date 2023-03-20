package config

type Cache struct {
	Redis     *Redis     `mapstructure:"redis"`
	Memcached *Memcached `mapstructure:"memcached"`
}

func (p Cache) Validate() error {
	options := enum{
		"redis":     p.Redis,
		"memcached": p.Memcached,
	}

	return options.Validate()
}
