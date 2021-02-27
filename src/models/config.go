package models

type Config struct {
	Server   ConfigServer   `yaml:"server"`
	Database ConfigDatabase `yaml:"database"`
}

type ConfigServer struct {
	Host string `yaml:"host" envconfig:"SERVER_HOST"`
	Port int    `yaml:"port" envconfig:"SERVER_PORT"`
}

type ConfigDatabase struct {
	Host string `yaml:"host" envconfig:"DATABASE_HOST"`
	User string `yaml:"user" envconfig:"DATABASE_USER"`
	Pass string `yaml:"pass" envconfig:"DATABASE_PASS"`
	Port int    `yaml:"port" envconfig:"DATABASE_PORT"`
	DB   string `yaml:"db" envconfig:"DATABASE_DB"`
}
