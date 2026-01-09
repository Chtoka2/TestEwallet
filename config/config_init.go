package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Env string `yaml:"env" env-required:"true" env-default:"local"`
	HttpServer `yaml:"http_server"`
}
type HttpServer struct{
	Address string `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func Config_init() *Config{
	var cfg Config
	if err := cleanenv.ReadConfig("config/local.yaml", &cfg); err != nil{
		log.Fatal("Cannot read config")
	}
	return &cfg
}