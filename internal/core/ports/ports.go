package ports

import (
// "marketflow/internal/core/domain"
)

type ExchangeRepository interface {
	GetData()
}

type ExchangeService interface {
	GetData()
}

type RedisRepository interface {
	Write(exchange string)
	Read() string
	LLen() int
}
