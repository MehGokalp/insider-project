package config

import (
	"github.com/rotisserie/eris"
	"os"
	"strconv"
)

func New() *Config {
	// create new Config from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		panic(eris.Wrap(err, "failed to parse PORT"))
	}

	return &Config{
		Mysql: Mysql{
			DSN: os.Getenv("MYSQL_DSN"),
		},
		Redis: Redis{
			DSN: os.Getenv("REDIS_DSN"),
		},
		Port: parsedPort,
		Env:  os.Getenv("ENV"),
		MessageProvider: MessageProvider{
			BaseUrl: os.Getenv("MESSAGE_PROVIDER_BASE_URL"),
		},
	}
}
