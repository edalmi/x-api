package config

import "errors"

type Prometheus struct{}

func (p Prometheus) Validate() error {
	return errors.New("not implemented")
}
