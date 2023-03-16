package config

type Logger struct {
	Std  *StdLogger `mapstructure:"std"`
	Zap  *Zap       `mapstructure:"zap"`
	Slog *Slog      `mapstructure:"slog"`
}

func (l Logger) Validate() error {
	providers := validator{
		"std":  l.Std,
		"zap":  l.Zap,
		"slog": l.Slog,
	}

	return providers.Validate()
}

type StdLogger struct{}

type Zap struct{}

type Slog struct{}
