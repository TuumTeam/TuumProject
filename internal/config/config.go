package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Server struct {
		Port         string `yaml:"port"`
		ReadTimeout  string `yaml:"read_timeout"`
		WriteTimeout string `yaml:"write_timeout"`
		IdleTimeout  string `yaml:"idle_timeout"`
	} `yaml:"server"`
	Database struct {
		Type            string `yaml:"type"`
		Path            string `yaml:"path"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	} `yaml:"database"`
	Auth struct {
		JWTSecret   string `yaml:"jwt_secret"`
		TokenExpiry string `yaml:"token_expiry"`
	} `yaml:"auth"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
	Security struct {
		CSRFProtection bool     `yaml:"csrf_protection"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"security"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal configuration data: %v", err)
	}
}
