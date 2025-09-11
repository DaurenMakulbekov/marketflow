package ports

import (
	"time"

	"marketflow/internal/core/domain"
)

type ExchangeRepository interface {
	GetFromExchange(exchange string) <-chan string
	Generator(exchange string)
	Close()
	CloseTest()
	GetExchanges() ([]string, []string)
	GetExchangesTest() ([]string, []string)
	GetExchangesBySymbol(symbol string) []string
	CheckSymbol(symbol string) bool
	CheckExchange(exchange string) bool
	CheckConnection() []string
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

	GetLowestPrice(symbol string) (domain.PriceSymbol, error)
	GetLowestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetLowestPriceByPeriod(symbol, period string) (domain.PriceSymbol, error)
	GetLowestExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error)

	GetAveragePrice(symbol string) (domain.PriceSymbol, error)
	GetAverageExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetAverageExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error)

	GetSystemStatus() domain.SystemStatus
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

	GetLowestPrice(symbol string) (domain.PriceSymbol, error)
	GetLowestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetLowestPriceByPeriod(symbol string, period time.Time) (domain.PriceSymbol, error)
	GetLowestExchangePriceByPeriod(exchange, symbol string, period time.Time) (domain.PriceExchangeSymbol, error)

	GetAveragePrice(symbol string) (domain.PriceSymbol, error)
	GetAverageExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error)
	GetAverageExchangePriceByPeriod(exchange, symbol string, period time.Time) (domain.PriceExchangeSymbol, error)
}

type Storage interface {
	Write(exchange domain.Exchange)
	GetAll() []domain.Exchange
	DeleteAll(exchanges []domain.Exchange)
	GetByPeriod(exchanges []string, symbol string, period int64) []domain.Exchange
	GetByExchangePeriod(exchange string, symbol string, period int64) []domain.Exchange
	GetLatest(exchanges []string, symbol string) []domain.Exchange
	GetLatestByExchange(exchange string, symbol string) []domain.Exchange
}
