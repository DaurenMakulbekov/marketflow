package storage

import (
	"marketflow/internal/core/domain"
)

type storage struct {
	table map[string]map[string]domain.Exchange
}

func NewStorage() *storage {
	var table = make(map[string]map[string]domain.Exchange)
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}
	var symbols = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	for i := range exchanges {
		for j := range symbols {
			table[exchanges[i]+":"+symbols[j]] = make(map[string]domain.Exchange)
		}
	}

	return &storage{
		table: table,
	}
}

func (st *storage) Write(exchange domain.Exchange) {
	st.table[exchange.Exchange+":"+exchange.Symbol][exchange.ID] = exchange
}

func (st *storage) GetAll() []domain.Exchange {
	var result []domain.Exchange

	for key := range st.table {
		for _, v := range st.table[key] {
			result = append(result, v)
		}
	}

	return result
}

func (st *storage) DeleteAll(exchanges []domain.Exchange) {
	for i := range exchanges {
		delete(st.table[exchanges[i].Exchange+":"+exchanges[i].Symbol], exchanges[i].ID)
	}
}
