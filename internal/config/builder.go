package config

import (
	"strings"

	"github.com/spf13/viper"
)

// db represents a database configuration.
type DB struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Name     string
}

// logger represents a logger configuration.
type logger struct {
	Level    string
	Encoding string
}

// Config is the configuration for the application.
type Config struct {
	Logger logger
	DB     DB
}

// New returns a new Config.
func New() *Config {
	confer := viper.New()

	confer.AutomaticEnv()
	confer.SetEnvPrefix("entain")
	confer.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := &Config{
		Logger: logger{
			Level:    confer.GetString("log.level"),
			Encoding: confer.GetString("log.encoding"),
		},
		DB: DB{
			Host:     confer.GetString("db.host"),
			Port:     confer.GetUint16("db.port"),
			User:     confer.GetString("db.user"),
			Password: confer.GetString("db.password"),
			Name:     confer.GetString("db.name"),
		},
	}

	return config
}
