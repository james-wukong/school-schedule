// Package config provides configuration loading and management for the application
// Package config provides functionality to load application configuration
package config

import (
	"github.com/james-wukong/school-schedule/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"databases"`
	Caches   CacheConfig    `mapstructure:"caches"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OTP      OtpConfig      `mapstructure:"otp"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Port        int    `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
	APIKey      string `mapstructure:"apikey"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Driver   string `mapstructure:"driver"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"` // String to preserve leading zeros
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	SSL      string `mapstructure:"ssl"`
	DSN      string `mapstructure:"dsn"`
	URL      string `mapstructure:"url"`
}

type CacheConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Expires  int    `mapstructure:"expires"`
}

type JWTConfig struct {
	Secret  string `mapstructure:"secret"`
	Expires int    `mapstructure:"expires"`
	Issuer  string `mapstructure:"issuer"`
	Refresh int    `mapstructure:"refresh"`
}

type OtpConfig struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func InitConfig() *Config {
	viper.SetConfigName("config") // Name of your file (config.yaml)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".") // Look in the current directory (project folder)

	// Enable environment variable overrides
	// Example: export APP_PORT=9000 will override the YAML
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Initialize logger for config loading with console output only
	conLog := logger.New(logger.LogConfig{
		EnableConsole: true,
		FilePath:      "",
		// MaxSize:       5,    // Rotate every 5MB
		// MaxBackups:    10,   // Keep last 10 files
		// Compress:      false, // Save disk space
	})

	if err := viper.ReadInConfig(); err != nil {
		conLog.Error().Err(err).Msg("Failed to read config file")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		conLog.Error().Err(err).Msg("Failed to unmarshal config")
	}

	return &config
}
