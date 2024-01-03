package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"prod"`
	StoragePath string `yaml:"storage_path" env-reauired:"true"`
	HTTPServer  `yaml:"http_server"`
	Clients     ClientConfig `yaml:"clients"`
	AppId       int64        `yaml:"appId" env-default:"1"`
	AppName     string       `yaml:"appName" env-default:"url-shortener"`
	AppSecret   string       `yaml:"appSecret" env-reauired:"true" env:"APP_SECRET"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	Insecure     bool          `yaml:"insecure"`
}

type ClientConfig struct {
	SSO Client `yaml:"sso"`
}

func MustLoad() Config {
	configPath := "D:/go_path/src/url-shortener/config/local.yaml"
	err := os.Setenv("CONFIG_PATH", configPath)
	if err != nil {
		log.Fatal("Failed to set CONFIG_PATH")
	}

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG PATH is not set")
	}

	// checl if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
