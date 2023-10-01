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
	Tracer
	Meter
	OtelExporter string
	Gmail
	Ses
}

type Logger struct {
	ServerMode string `mapstructure:"SERVER_MODE"`
	Encoding   string `mapstructure:"LOG_ENCODING"`
	Level      string `mapstructure:"LOG_LEVEL"`
}

type Tracer struct {
	Name string `mapstructure:"TRACE_NAME"`
}

type Meter struct {
	Name string `mapstructure:"METER_NAME"`
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

	log.Println("env variable DEV=false reading from system")
	return loadFromEnv()
}

func loadFromFile(path string) *Config {

	var config Config

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("LoadFile.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("LoadFile.Unmarshal: %v", err)
	}

	return &config
}

func loadFromEnv() *Config {

	return &Config{}
}
