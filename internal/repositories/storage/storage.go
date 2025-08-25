package storage

import (
	"marketflow/internal/core/domain"
)

type storage struct {
	table map[string]map[string][]domain.Exchange
}

func NewStorage() *storage {
	var table = make(map[string]map[string][]domain.Exchange)
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}

	for i := range exchanges {
		table[exchanges[i]] = make(map[string][]domain.Exchange)
	}

	return &storage{
		table: table,
	}
}

func (st *storage) Write(exchange domain.Exchange) {
	st.table[exchange.Exchange][exchange.Symbol] = append(st.table[exchange.Exchange][exchange.Symbol], exchange)
}

func (st *storage) GetAll(exchanges []string) []domain.Exchange {
	var result []domain.Exchange

	for i := range exchanges {
		for k := range st.table[exchanges[i]] {
			result = append(result, st.table[exchanges[i]][k]...)
		}
	}

	return result
}
