package config

import "errors"

type Prometheus struct {
	Namespace string `mapstructure:"namespace"`
}

func (p Prometheus) Validate() error {
	return errors.New("not implemented")
}
