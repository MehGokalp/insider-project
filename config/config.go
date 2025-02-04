package config

type Config struct {
	Mysql           Mysql
	Redis           Redis
	Port            int
	Env             string
	MessageProvider MessageProvider
}

type MessageProvider struct {
	BaseUrl string
}

type Mysql struct {
	DSN string
}

type Redis struct {
	DSN string
}
