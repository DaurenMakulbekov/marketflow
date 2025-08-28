package ports

import (
	"marketflow/internal/core/domain"
)

type ExchangeRepository interface {
	GetFromExchange(exchange string) <-chan string
	Generator() <-chan string
}

type ExchangeService interface {
	LiveMode()
	TestMode()
}

type RedisRepository interface {
	Write(exchange domain.Exchange) error
	ReadAll(exchanges, pairName []string) ([]domain.Exchange, error)
	DeleteAll(exchanges []domain.Exchange) error
	CheckConnection() error
	Reconnect()
}

type PostgresRepository interface {
	Write(exchange []domain.Exchanges) error
}

type Storage interface {
	Write(exchange domain.Exchange)
	GetAll() []domain.Exchange
	DeleteAll(exchanges []domain.Exchange)
}
