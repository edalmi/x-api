package config

type Cache struct {
	Redis     *Redis     `mapstructure:"redis"`
	Memcached *Memcached `mapstructure:"memcached"`
}

func (p Cache) Validate() error {
	parts := validator{
		"redis":     p.Redis,
		"memcached": p.Memcached,
	}

	return parts.Validate()
}
