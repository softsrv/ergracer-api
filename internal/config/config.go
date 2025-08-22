package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Mailgun  MailgunConfig  `yaml:"mailgun"`
	App      AppConfig      `yaml:"app"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type MailgunConfig struct {
	Domain    string `yaml:"domain"`
	APIKey    string `yaml:"api_key"`
	FromEmail string `yaml:"from_email"`
	FromName  string `yaml:"from_name"`
}

type AppConfig struct {
	URL  string `yaml:"url"`
	Port int    `yaml:"port"`
}

// Legacy getters for backward compatibility
func (c *Config) DatabaseURL() string {
	return c.Database.URL
}

func (c *Config) JWTSecret() string {
	return c.JWT.Secret
}

func (c *Config) MailgunDomain() string {
	return c.Mailgun.Domain
}

func (c *Config) MailgunAPIKey() string {
	return c.Mailgun.APIKey
}

func (c *Config) MailgunFromEmail() string {
	return c.Mailgun.FromEmail
}

func (c *Config) MailgunFromName() string {
	return c.Mailgun.FromName
}

func (c *Config) AppURL() string {
	return c.App.URL
}

func Load() *Config {
	configPath := getConfigPath()
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file '%s': %v", configPath, err))
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(fmt.Sprintf("Failed to parse config file '%s': %v", configPath, err))
	}
	
	return &config
}

func getConfigPath() string {
	if configPath := os.Getenv("CONFIG_FILE_PATH"); configPath != "" {
		return configPath
	}
	return "config.yaml"
}