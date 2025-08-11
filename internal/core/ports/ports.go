package ports

import (
	"marketflow/internal/core/domain"
)

type ExchangeRepository interface {
	LiveMode()
	TestMode()
}

type ExchangeService interface {
	LiveMode()
	TestMode()
}

type RedisRepository interface {
	Write(exchange domain.Exchange) error
	ReadAll(exchanges, pairName []string) ([]domain.Exchange, error)
	DeleteAll(exchanges []domain.Exchange) error
}

type PostgresRepository interface {
	Write(exchange []domain.Exchanges) error
}
