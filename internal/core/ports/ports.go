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
	Write(exchange string)
	Read() string
	LLen() int
}

type PostgresRepository interface {
	Write(exchange domain.Exchanges) error
}
