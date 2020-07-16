package config

// Provider is foobar
type Provider interface {
	Unmarshal(target interface{}) error
}

// Config is foobar
type Config struct {
	Providers []Provider
}

// Unmarshal is foobar
func (c Config) Unmarshal(target interface{}) error {
	for _, provider := range c.Providers {
		if err := provider.Unmarshal(target); err != nil {
			return err
		}
	}

	return nil
}
