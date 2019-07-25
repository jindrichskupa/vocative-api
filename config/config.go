package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config Application configuration structure
type Config struct {
	DB   *DBConfig
	IP   string
	Port uint16
}

// DBConfig structure
type DBConfig struct {
	Dialect  string
	Hostname string
	Port     uint16
	Username string
	Password string
	Name     string
	Charset  string
	Prefix   string
	Retries  int
}

// EnvConfig stores config from ENV
type EnvConfig struct {
	DBHostname string `envconfig:"db_hostname" required:"false" default:"localhost"`
	DBPort     uint16 `envconfig:"db_port" required:"false" default:"5432"`
	DBName     string `envconfig:"db_name" required:"false" default:"vocativedb"`
	DBUser     string `envconfig:"db_user" required:"false" default:"postgres"`
	DBPassword string `envconfig:"db_password" required:"false" default:"password"`
	DBRetries  int    `envconfig:"db_retries" required:"false" default:"1"`
	ListenIP   string `envconfig:"listen_ip" required:"false" default:"0.0.0.0"`
	ListenPort uint16 `envconfig:"listen_port" required:"false" default:"8080"`
}

// GetConfig get Application configuration
func GetConfig() *Config {
	var s EnvConfig
	err := envconfig.Process("vocative", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	config := Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Hostname: s.DBHostname,
			Port:     s.DBPort,
			Username: s.DBUser,
			Password: s.DBPassword,
			Name:     s.DBName,
			Retries:  s.DBRetries,
			Charset:  "utf8",
			Prefix:   "",
		},
		IP:   s.ListenIP,
		Port: s.ListenPort,
	}

	return &config
}

// ListenAddress returns listen string
func (c *Config) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}
