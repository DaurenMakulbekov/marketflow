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
}

type ExchangeService interface {
	LiveMode()
	TestMode()

	GetLatestSymbol(symbol string) (domain.Exchange, error)
	GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error)
	GetHighestSymbol(symbol string) (domain.PriceSymbol, error)
	GetHighestExchangeSymbol(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetHighestSymbolByPeriod(symbol, period string) (domain.PriceSymbol, error)
	GetHighestExchangeSymbolByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error)
}

type RedisRepository interface {
	Write(exchange domain.Exchange) error
	ReadAll(exchanges, pairName []string) ([]domain.Exchange, error)
	DeleteAll(exchanges []domain.Exchange) error
	CheckConnection() error
	Reconnect()

	GetLatestSymbol(exchanges []string, symbol string) ([]domain.Exchange, error)
	GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error)
	GetPriceByPeriod(exchanges []string, symbol, period string) ([]domain.Exchange, error)
}

type PostgresRepository interface {
	Write(exchange []domain.Exchanges) error
	GetHighestSymbol(symbol string) (domain.PriceSymbol, error)
	GetHighestExchangeSymbol(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetHighestSymbolByPeriod(symbol string, period time.Time) (domain.PriceSymbol, error)
	GetHighestExchangeSymbolByPeriod(exchange, symbol string, period time.Time) (domain.PriceExchangeSymbol, error)
}

type Storage interface {
	Write(exchange domain.Exchange)
	GetAll() []domain.Exchange
	DeleteAll(exchanges []domain.Exchange)
}
