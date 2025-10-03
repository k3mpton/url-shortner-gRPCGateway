package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	GRPC     Grpc     `yaml:"gRpc" env-required:"true"`
	Database Database `yaml:"database" env-required:"true"`
}

type Database struct {
	URL string `env:"DATABASE_URL"`
}

type Grpc struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustReadCfg() *Config {
	path := fetchPathCfg()

	if _, err := os.Stat(path); err != nil {
		panic(err)
	}

	cfg := Config{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

var (
	pathCfg = flag.String("cfg", "./config/config.yaml", "config path")
)

func fetchPathCfg() string {
	flag.Parse()
	return *pathCfg
}
