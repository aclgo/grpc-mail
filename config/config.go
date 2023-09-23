package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceHTTPPort int           `mapstructure:"API_HTTP_PORT"`
	ServiceGRPCPort int           `mapstructure:"API_GRPC_PORT"`
	IntervalSend    time.Duration `mapstructure:"INTERVAL_SEND"`
	Logger
	Gmail
	Ses
}

type Logger struct {
	ServerMode string `mapstructure:"SERVER_MODE"`
	Encoding   string `mapstructure:"LOG_ENCODING"`
	Level      string `mapstructure:"LOG_LEVEL"`
}

type Redis struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

type Ses struct {
}

type Gmail struct {
	Identity string `mapstructure:"GMAIL_IDENTITY"`
	Username string `mapstructure:"GMAIL_USERNAME"`
	Password string `mapstructure:"GMAIL_PASSWORD"`
	Host     string `mapstructure:"GMAIL_HOST"`
	Port     int    `mapstructure:"GMAIL_PORT"`
}

func Load(path string) *Config {

	if os.Getenv("DEV") == "true" {
		return loadFromFile(path)
	}

	return loadFromEnv()
}

func loadFromFile(path string) *Config {

	var config Config

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("LoadFile.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("LoadFile.Unmarshal: %v", err)
	}

	return nil
}

func loadFromEnv() *Config {
	return &Config{}
}
