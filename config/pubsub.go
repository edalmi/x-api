package config

type Pubsub struct {
	Redis    *Redis    `mapstructure:"redis"`
	RabbitMQ *RabbitMQ `mapstructure:"rabbitmq"`
}

func (p Pubsub) Validate() error {
	parts := validator{
		"redis":    p.Redis,
		"rabbitmq": p.RabbitMQ,
	}

	return parts.Validate()
}
