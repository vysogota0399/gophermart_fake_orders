package main

import (
	"github.com/vysogota0399/gophermart/internal/config"
	"github.com/vysogota0399/gophermart/internal/seeds"
	"github.com/vysogota0399/gophermart/internal/storage"
)

func main() {
	cfg := config.MustNewConfig()
	strg, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	s := seeds.ProductsCreator{Strg: strg, AccrualCreateGoodsURL: cfg.AccrualCreateGoodsURL}
	if err := s.Call(); err != nil {
		panic(err)
	}
}
