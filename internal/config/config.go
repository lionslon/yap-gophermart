package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"RUN_ADDRESS"`
	Accrual string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DSN     string `env:"DATABASE_URI"`
	Key     []byte `env:"KEY"`
}

func GetConfig() *Config {

	cfg := &Config{}

	var key string
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Gophermart service address and port")
	flag.StringVar(&cfg.Accrual, "r", "localhost:8078", "Accrual service address and port")
	flag.StringVar(&cfg.DSN, "d", "''", "Postgresql DSN string")
	flag.StringVar(&key, "k", "gophermart", "Secret key")
	cfg.Key = []byte(key)

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		fmt.Println(err)
	}

	return cfg

}
