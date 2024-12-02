package config

type Config struct {
	host        `mapstructure:",squash"`
	Development bool
	PostgreSQL  postgreSQL
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
