package ports

import (
	"marketflow/internal/core/domain"
)

type ExchangeRepository interface {
	GetFromExchange(exchange string) <-chan string
	Generator() <-chan string
	CloseTest()
}

type ExchangeService interface {
	LiveMode()
	TestMode()

	GetLatestSymbol(symbol string) (domain.Exchange, error)
	GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error)
	GetHighestSymbol(symbol string) (domain.Price, error)
}

type RedisRepository interface {
	Write(exchange domain.Exchange) error
	ReadAll(exchanges, pairName []string) ([]domain.Exchange, error)
	DeleteAll(exchanges []domain.Exchange) error
	CheckConnection() error
	Reconnect()

	GetLatestSymbol(exchange, symbol string) (domain.Exchange, error)
}

type PostgresRepository interface {
	Write(exchange []domain.Exchanges) error
	GetHighestSymbol(symbol string) (domain.Price, error)
}

type Storage interface {
	Write(exchange domain.Exchange)
	GetAll() []domain.Exchange
	DeleteAll(exchanges []domain.Exchange)
}
