package config

type Queue struct {
	Redis    *Redis    `mapstructure:"redis"`
	RabbitMQ *RabbitMQ `mapstructure:"rabbitmq"`
}

func (l Queue) Validate() error {
	providers := enum{
		"redis":    l.Redis,
		"rabbitmq": l.RabbitMQ,
	}

	return providers.Validate()
}
