package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`

	// Logging configuration
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	LogFormat      string `mapstructure:"LOG_FORMAT"`
	LogOutput      string `mapstructure:"LOG_OUTPUT"`
	LogRequestBody bool   `mapstructure:"LOG_REQUEST_BODY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values for logging
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_OUTPUT", "stdout")
	viper.SetDefault("LOG_REQUEST_BODY", false)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
