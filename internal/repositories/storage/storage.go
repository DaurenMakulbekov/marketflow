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
	var exchangesTest = []string{"exchange1_test", "exchange2_test", "exchange3_test"}
	var symbolsTest = []string{"BTCUSDT_test", "DOGEUSDT_test", "TONUSDT_test", "SOLUSDT_test", "ETHUSDT_test"}

	for i := range exchanges {
		for j := range symbols {
			table[exchanges[i]+":"+symbols[j]] = make(map[string]domain.Exchange)
			table[exchangesTest[i]+":"+symbolsTest[j]] = make(map[string]domain.Exchange)
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

func (st *storage) GetByPeriod(exchanges []string, symbol string, period int64) []domain.Exchange {
	var result []domain.Exchange

	for key := range st.table {
		for _, v := range st.table[key] {
			for i := range exchanges {
				if exchanges[i] == v.Exchange && symbol == v.Symbol && v.Timestamp >= period {
					result = append(result, v)
				}
			}
		}
	}

	return result
}

func (st *storage) GetByExchangePeriod(exchange string, symbol string, period int64) []domain.Exchange {
	var result []domain.Exchange

	for key := range st.table {
		for _, v := range st.table[key] {
			if exchange == v.Exchange && symbol == v.Symbol && v.Timestamp >= period {
				result = append(result, v)
			}
		}
	}

	return result
}

func (st *storage) GetLatest(exchanges []string, symbol string) []domain.Exchange {
	var result []domain.Exchange

	for key := range st.table {
		for _, v := range st.table[key] {
			for i := range exchanges {
				if exchanges[i] == v.Exchange && symbol == v.Symbol {
					result = append(result, v)
				}
			}
		}
	}

	return result
}

func (st *storage) GetLatestByExchange(exchange string, symbol string) []domain.Exchange {
	var result []domain.Exchange

	for key := range st.table {
		for _, v := range st.table[key] {
			if exchange == v.Exchange && symbol == v.Symbol {
				result = append(result, v)
			}
		}
	}

	return result
}

func (st *storage) DeleteAll(exchanges []domain.Exchange) {
	for i := range exchanges {
		delete(st.table[exchanges[i].Exchange+":"+exchanges[i].Symbol], exchanges[i].ID)
	}
}
