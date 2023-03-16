package config

import "fmt"

type Pubsub struct {
	Redis    *Redis      `mapstructure:"redis"`
	RabbitMQ interface{} `mapstructure:"rabbitmq"`
}

func (p Pubsub) Validate() error {
	parts := validator{
		"redis":    p.Redis,
		"rabbitmq": p.RabbitMQ,
	}

	return parts.Validate()
}

type validator map[string]interface{}

func (v validator) Validate() error {
	var defined bool
	var selected interface{}
	for k, v := range v {
		if v != nil {
			if defined {
				return fmt.Errorf("%v only one provider should be configured", k)
			}

			defined = true
			selected = v
		}
	}

	if selected != nil {
		if s, ok := selected.(interface{ Validate() error }); ok {
			return s.Validate()
		}
	}

	return nil
}
