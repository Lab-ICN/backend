package config

type Config struct {
	host
	PostgreSQL postgreSQL
	Logging    logging
}

type host struct {
	Address string
	Port    int
}

type postgreSQL struct {
	host
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
