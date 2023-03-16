package config

import "fmt"

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
