package config

import (
	"os"
)

type DB struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type Redis struct {
	Addr     string
	Password string
}

type Exchange struct {
	Name string
	Host string
	Port string
}

type AppConfig struct {
	DB        *DB
	Redis     *Redis
	Exchanges []*Exchange
}

func NewDB() *DB {
	var db = &DB{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	return db
}

func NewRedis() *Redis {
	var rdb = &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	return rdb
}

func NewExchanges() []*Exchange {
	var exchanges []*Exchange

	var exchange1 = &Exchange{
		Name: os.Getenv("EXCHANGE1_NAME"),
		Host: os.Getenv("EXCHANGE1_HOST"),
		Port: os.Getenv("EXCHANGE1_PORT"),
	}

	var exchange2 = &Exchange{
		Name: os.Getenv("EXCHANGE2_NAME"),
		Host: os.Getenv("EXCHANGE2_HOST"),
		Port: os.Getenv("EXCHANGE2_PORT"),
	}

	var exchange3 = &Exchange{
		Name: os.Getenv("EXCHANGE3_NAME"),
		Host: os.Getenv("EXCHANGE3_HOST"),
		Port: os.Getenv("EXCHANGE3_PORT"),
	}

	exchanges = append(exchanges, exchange1, exchange2, exchange3)

	return exchanges
}

func NewAppConfig() *AppConfig {
	var config = &AppConfig{
		DB:        NewDB(),
		Redis:     NewRedis(),
		Exchanges: NewExchanges(),
	}

	return config
}
