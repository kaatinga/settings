package env_loader

import "errors"

type Settings struct {
	strings map[string]string
	fixed   bool
}

func (s Settings) GetString(setting string) (value string, err error) {
	// TBI
	return
}

func (s Settings) GetByte(setting string) (value byte, err error) {
	// TBI
	return
}

// SetString sets an setting that must be in the environment variables
func (s Settings) SetString(value string) error {
	switch s.fixed {
	case true:
		return errors.New("impossible to add a setting as settings are read already")
	case false:
		// TBI
	}
	return nil
}
