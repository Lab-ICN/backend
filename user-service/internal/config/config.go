package config

type Config struct {
	PostgreSQL  postgreSQL
	JwtKey      string
	ApiKey      string
	host        `mapstructure:",squash"`
	Development bool
}

type host struct {
	Address string
	Port    int
}

type postgreSQL struct {
	Username string
	Password string
	Database string
	host     `mapstructure:",squash"`
}
