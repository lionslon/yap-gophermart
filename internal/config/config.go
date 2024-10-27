package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"time"
)

type Config struct {
	Address         string `env:"RUN_ADDRESS"`
	Accrual         string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DSN             string `env:"DATABASE_URI"`
	Key             []byte `env:"KEY"`
	AccrualInterval int
	TokenExp        time.Duration
}

func GetConfig() *Config {

	cfg := &Config{}

	var key string
	var tokExp int

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Gophermart service address and port")
	flag.StringVar(&cfg.Accrual, "r", "localhost:8078", "Accural service address and port")
	flag.StringVar(&cfg.DSN, "d", "", "Postgresql DSN string")
	flag.StringVar(&key, "k", "gophermart", "Secret key")
	flag.IntVar(&cfg.AccrualInterval, "i", 2, "This is timeout between requests to the accrual service")
	flag.IntVar(&tokExp, "t", 1, "This is jwt token expiration")

	cfg.Key = []byte(key)
	cfg.TokenExp = time.Hour * time.Duration(tokExp)

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		fmt.Println(err)
	}

	return cfg

}
