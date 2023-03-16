package config

type Queue struct {
	Redis    *Redis    `mapstructure:"redis"`
	RabbitMQ *RabbitMQ `mapstructure:"rabbitmq"`
}

func (l Queue) Validate() error {
	providers := validator{
		"redis":    l.Redis,
		"rabbitmq": l.RabbitMQ,
	}

	return providers.Validate()
}
