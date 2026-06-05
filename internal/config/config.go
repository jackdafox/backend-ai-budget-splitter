package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	App         AppConfig         `mapstructure:"app"`
	Database    DatabaseConfig    `mapstructure:"database"`
	AIProviders AIProvidersConfig `mapstructure:"ai_providers"`
	Auth        AuthConfig        `mapstructure:"auth"`
}

// AppConfig holds application-level configuration.
type AppConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// AIProvidersConfig holds AI provider configurations.
type AIProvidersConfig struct {
	OpenAI    ProviderConfig `mapstructure:"openai"`
	Anthropic ProviderConfig `mapstructure:"anthropic"`
}

// ProviderConfig holds a single AI provider configuration.
type ProviderConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	KeyExpireDays int    `mapstructure:"key_expire_days"`
}

// Load reads configuration from file and environment variables.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Enable environment variable override
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
