package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	// Interval defines the interval for the job to post process transactions.
	Interval   = 1
	SourceType = "Source-Type"
)

// db represents a database configuration.
type db struct {
	Host     string
	Port     string
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
	DB     db
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
		DB: db{
			Host:     confer.GetString("db.host"),
			Port:     confer.GetString("db.port"),
			User:     confer.GetString("db.user"),
			Password: confer.GetString("db.password"),
			Name:     confer.GetString("db.name"),
		},
	}

	return config
}
