package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	GitHub   GitHubConfig   `mapstructure:"github"`
}

type ServerConfig struct {
	Port           int      `mapstructure:"port"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type GitHubConfig struct {
	Token      string `mapstructure:"token"`
	Owner      string `mapstructure:"owner"`
	Prefix     string `mapstructure:"prefix"`
	GitOpsRepo string `mapstructure:"gitops_repo"`
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./cmd/api")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "lokiforce")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "super_secret_signing_key_for_lokiforce")
	viper.SetDefault("github.token", "")
	viper.SetDefault("github.owner", "lokiforce-dev")
	viper.SetDefault("github.prefix", "lkf-")
	viper.SetDefault("github.gitops_repo", "gitops-infra")
	viper.SetDefault("server.allowed_origins", []string{"http://localhost:5173"})

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: failed to read config file: %v. Using defaults.", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
