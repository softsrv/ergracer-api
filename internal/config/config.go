package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	App      AppConfig      `yaml:"app"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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

func (c *Config) SMTPHost() string {
	return c.SMTP.Host
}

func (c *Config) SMTPPort() string {
	return fmt.Sprintf("%d", c.SMTP.Port)
}

func (c *Config) SMTPUsername() string {
	return c.SMTP.Username
}

func (c *Config) SMTPPassword() string {
	return c.SMTP.Password
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