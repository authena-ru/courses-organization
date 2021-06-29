package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort      = "8080"
	defaultHTTPRWTimeout = 10 * time.Second
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
	}

	MongoConfig struct {
		URI          string
		Username     string
		Password     string
		DatabaseName string
	}

	HTTPConfig struct {
		Host         string
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
)

func New(configsDir string) (*Config, error) {
	setDefaults()

	if err := parseEnv(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, viper.GetString("environment")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshall(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}

func parseEnv() error {
	if err := parseMongoFromEnv(); err != nil {
		return err
	}

	return parseAppFromEnv()
}

func parseMongoFromEnv() error {
	viper.SetEnvPrefix("mongo")

	if err := viper.BindEnv("uri"); err != nil {
		return err
	}

	if err := viper.BindEnv("username"); err != nil {
		return err
	}

	return viper.BindEnv("password")
}

func parseAppFromEnv() error {
	viper.SetEnvPrefix("app")

	return viper.BindEnv("environment")
}

func parseConfigFile(configsDir, env string) error {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName(env)

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.Environment = viper.GetString("environment")
	cfg.Mongo.URI = viper.GetString("uri")
	cfg.Mongo.Username = viper.GetString("username")
	cfg.Mongo.Password = viper.GetString("password")
}

func unmarshall(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	return viper.UnmarshalKey("mongo", &cfg.Mongo)
}
