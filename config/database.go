package config

type DB struct {
	Postgres *Postgres `mapstructure:"postgres"`
	MySQL    *MySQL    `mapstructure:"mysql"`
	SQLite   *SQLite   `mapstructure:"sqlite"`
	MariaDB  *MariaDB  `mapstructure:"mariadb"`
}

func (p DB) Validate() error {
	parts := enum{
		"postgres": p.Postgres,
		"mysql":    p.MySQL,
		"sqlite":   p.SQLite,
		"mariadb":  p.MariaDB,
	}

	return parts.Validate()
}
