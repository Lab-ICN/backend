package config

type Config struct {
	Key            string
	GoogleClientID string
	PostgreSQL     postgreSQL
	LogPath        string
	host           `mapstructure:",squash"`
	JWT            jwt
	Development    bool
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

type jwt struct {
	Key        string
	AccessTTL  int
	RefreshTTL int
}
