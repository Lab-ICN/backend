package config

type Config struct {
	host       `mapstructure:",squash"`
	PostgreSQL postgreSQL
	Logging    logging
}

type host struct {
	Address string
	Port    int
}

type postgreSQL struct {
	host     `mapstructure:",squash"`
	Username string
	Password string
	Database string
}

type logging struct {
	Level             string
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	OutputPaths       []string
	ErrorOutputPaths  []string
}
