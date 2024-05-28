package config

type Config struct {
	ConnString string
}

func NewConfig() (*Config, error) {
	return &Config{
		ConnString: "postgresql://mission:password@localhost:5432/mission_db?sslmode=disable",
	}, nil
}
