package config

type Memcached struct {
	Addresses []string `mapstructure:"addresses"`
}
