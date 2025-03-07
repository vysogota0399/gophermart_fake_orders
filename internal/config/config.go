package config

import "github.com/caarlos0/env"

type Config struct {
	GRPCAddress           string `json:"grpc_address" env:"GOPHERMART_FAKE_ORDERS_GRPC_ADDRESS" envDefault:"127.0.0.1:8050"`
	HTTPAddress           string `json:"http_address" env:"GOPHERMART_FAKE_ORDERS_HTTP_ADDRESS" envDefault:"127.0.0.1:8060"`
	AccrualCreateGoodsURL string `json:"accrual_create_goods_url" env:"GOPHERMART_FAKE_ORDERS_ACCRUAL_CREATE_GOODS_URL" envDefault:"http://127.0.0.1:8080/api/goods"`
	LogLevel              int   `json:"log_level" env:"GOPHERMART_FAKE_ORDERS_LOG_LEVEL" envDefault:"-1"`
	DatabaseDSN           string `json:"database_dsn" env:"GOPHERMART_FAKE_ORDERS_DATABASE_DSN" envDefault:"postgres://postgres:secret@127.0.0.1:5432/gophermart_development"`
}

func MustNewConfig() *Config {
	c := &Config{}
	env.Parse(c)

	return c
}
