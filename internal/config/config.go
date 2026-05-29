package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env"`
	HttpServer HttpServer `yaml:"http_server"`
	DB         PrimaryDB  `yaml:"database"`
	MINIOS3    Minio      `yaml:"minio_s3"`
	Clients    Clients    `yaml:"clients"`
}
type Clients struct {
	Auth Client `yaml:"auth_client"`
}
type Client struct {
	Address  string        `yaml:"address"`
	Timeout  time.Duration `yaml:"timeout" env-default:"5s"`
	Retries  int           `yaml:"retries" env-default:"5"`
	Insecure bool          `yaml:"insecure" env-default:"true"`
}
type HttpServer struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"5000"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}
type Minio struct {
	URL        string `yaml:"url"`
	ACCESS_KEY string `yaml:"access_key"`
	SECRET_KEY string `yaml:"secret_key"`
	SSL        bool   `yaml:"ssl"`
}
type PrimaryDB struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"dbname"`
	SSL         string `yaml:"ssl"`
	MaxAttempts int    `yaml:"max_attempts"`
}

func InitConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env.dev"
	}
	fmt.Println("env name", envFile)
	if err := godotenv.Load(envFile); err != nil {
		if !os.IsNotExist(err) {
			slog.Error("ошибка при инициализации переменных окружения", slog.Any("err", err))
		}
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}

	// Read YAML file and substitute ${VAR} with environment variables
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg
}
