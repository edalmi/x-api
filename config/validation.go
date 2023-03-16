package config

import "fmt"

type enum map[string]interface{}

func (e enum) Validate() error {
	var alreadyDefined bool

	for k, option := range e {
		if option != nil {
			if alreadyDefined {
				fmt.Printf("%#v", option)
				return fmt.Errorf("%v only one provider should be configured", k)
			}

			if s, ok := option.(interface{ Validate() error }); ok {
				return s.Validate()
			}

			alreadyDefined = true
		}
	}

	return nil
}
