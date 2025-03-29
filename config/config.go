package config

import (
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/spf13/viper"
)

type Config struct {
	// Database Configuration
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	// Token Configuration
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	// Server Configuration
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	Environment       string `mapstructure:"ENVIROMENT"`

	// Blockchain Configuration
	CryptDeployURL     string `mapstructure:"CRYPT_DEPLOY_URL"`
	ContractPrivateKey string `mapstructure:"CONTRACT_PRIVATE_KEY"`
	ContractAddress    string `mapstructure:"CONTRACT_ADDRESS"`
	DeployAddr         string `mapstructure:"DEPLOY_ADDR"`
	DeployPrivateKey   string `mapstructure:"DEPLOY_PRIVATE_KEY"`
	PassPhase          string `mapstructure:"PASS_PHASE"`

	// Storage Configuration
	CloudinaryAPIKey string `mapstructure:"CLOUDINARY_API_KEY"`

	// Email Configuration
	Email     string `mapstructure:"EMAIL"`
	EmailPass string `mapstructure:"EMAIL_PASS"`
	SMTPPort  int `mapstructure:"SMTP_PORT"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPName  string `mapstructure:"SMTP_NAME"`
	SenderEmail  string `mapstructure:"SENDER_EMAIL"`
	SenderName  string `mapstructure:"SENDER_NAME"`
	SMTPUsername  string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`

	// Security Configuration
	RotateRefreshTokens bool `mapstructure:"ROTATE_REFRESH_TOKENS"`
	MaxLoginAttempts    int  `mapstructure:"MAX_LOGIN_ATTEMPTS"`

	// OTP Configuration
	OTPExpiryDuration time.Duration `mapstructure:"OTP_EXPIRY_DURATION"`
	MaxOTPAttempts    int           `mapstructure:"MAX_OTP_ATTEMPTS"`

	// Logging configuration
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	LogFormat      string `mapstructure:"LOG_FORMAT"`
	LogOutput      string `mapstructure:"LOG_OUTPUT"`
	LogRequestBody bool   `mapstructure:"LOG_REQUEST_BODY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")

	// Set the path to look for the configuration file
	viper.AddConfigPath(path)

	// Set default values
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgres://root:secret@localhost:5433/defi?sslmode=disable")
	viper.SetDefault("HTTP_SERVER_ADDRESS", ":8080")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "15m")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "24h")
	viper.SetDefault("ENVIROMENT", "development")
	viper.SetDefault("ROTATE_REFRESH_TOKENS", true)
	viper.SetDefault("MAX_LOGIN_ATTEMPTS", 5)
	viper.SetDefault("OTP_EXPIRY_DURATION", "5m")
	viper.SetDefault("MAX_OTP_ATTEMPTS", 3)

	// Set default values for logging
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_OUTPUT", "stdout")
	viper.SetDefault("LOG_REQUEST_BODY", false)
	// Enable VIPER to read environment variables
	viper.AutomaticEnv()

	// Set the type of the configuration file
	viper.SetConfigType("env")

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the configuration into the Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Parse duration strings
	config.AccessTokenDuration, err = time.ParseDuration(viper.GetString("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return
	}

	config.RefreshTokenDuration, err = time.ParseDuration(viper.GetString("REFRESH_TOKEN_DURATION"))
	if err != nil {
		return
	}

	config.OTPExpiryDuration, err = time.ParseDuration(viper.GetString("OTP_EXPIRY_DURATION"))
	if err != nil {
		return
	}

	return
}

// ToDomainConfig converts the config package's Config to domain.Config
func (c *Config) ToDomainConfig() *domain.Config {
	return &domain.Config{
		// JWT Configuration
		TokenSymmetricKey:    c.TokenSymmetricKey,
		AccessTokenDuration:  c.AccessTokenDuration,
		RefreshTokenDuration: c.RefreshTokenDuration,

		// Database Configuration
		DBDriver: c.DBDriver,
		DBSource: c.DBSource,

		// Server Configuration
		HTTPServerAddress: c.HTTPServerAddress,
		Environment:       c.Environment,

		// Security Configuration
		RotateRefreshTokens: c.RotateRefreshTokens,
		MaxLoginAttempts:    c.MaxLoginAttempts,

		// OTP Configuration
		OTPExpiryDuration: c.OTPExpiryDuration,
		MaxOTPAttempts:    c.MaxOTPAttempts,
	}
}
