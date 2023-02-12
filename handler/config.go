package handler

type Config struct {
	errorWriter ErrorWriter
}

type Option func(*Config)

func WithErrorWriter(ew ErrorWriter) Option {
	return func(c *Config) {
		c.errorWriter = ew
	}
}
