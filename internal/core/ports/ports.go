package ports

import (
	"marketflow/internal/core/domain"
	"time"
)

type ExchangeRepository interface {
	GetFromExchange(exchange string) <-chan string
	Generator() <-chan string
	CloseTest()
	GetExchangesBySymbol(symbol string) []string
	CheckSymbol(symbol string) bool
	CheckExchange(exchange string) bool
}

type ExchangeService interface {
	LiveMode()
	TestMode()

	GetLatestPrice(symbol string) (domain.Exchange, error)
	GetLatestExchangePrice(exchange, symbol string) (domain.Exchange, error)

	GetHighestPrice(symbol string) (domain.PriceSymbol, error)
	GetHighestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetHighestPriceByPeriod(symbol, period string) (domain.PriceSymbol, error)
	GetHighestExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error)
}

type RedisRepository interface {
	Write(exchange domain.Exchange) error
	ReadAll(exchanges, pairName []string) ([]domain.Exchange, error)
	DeleteAll(exchanges []domain.Exchange) error
	CheckConnection() error
	Reconnect()

	GetLatestPrice(exchanges []string, symbol string) ([]domain.Exchange, error)
	GetLatestExchangePrice(exchange, symbol string) (domain.Exchange, error)

	GetPriceByPeriod(exchanges []string, symbol, period string) ([]domain.Exchange, error)
	GetExchangePriceByPeriod(exchange, symbol, period string) ([]domain.Exchange, error)
}

type PostgresRepository interface {
	Write(exchange []domain.Exchanges) error
	GetHighestPrice(symbol string) (domain.PriceSymbol, error)
	GetHighestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetHighestPriceByPeriod(symbol string, period time.Time) (domain.PriceSymbol, error)
	GetHighestExchangePriceByPeriod(exchange, symbol string, period time.Time) (domain.PriceExchangeSymbol, error)
}

type Storage interface {
	Write(exchange domain.Exchange)
	GetAll() []domain.Exchange
	DeleteAll(exchanges []domain.Exchange)
}
