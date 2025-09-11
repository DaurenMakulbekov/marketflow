package storage

import (
	"sync"

	"marketflow/internal/core/domain"
	"marketflow/internal/infrastructure/config"
)

type storage struct {
	table map[string]map[string]domain.Exchange
	mu    sync.Mutex
}

func NewStorage(configs []*config.Exchange) *storage {
	table := make(map[string]map[string]domain.Exchange)
	var exchanges []string
	symbols := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}
	var exchangesTest []string
	symbolsTest := []string{"BTCUSDT_test", "DOGEUSDT_test", "TONUSDT_test", "SOLUSDT_test", "ETHUSDT_test"}

	for i := range configs {
		if i < 3 {
			exchanges = append(exchanges, configs[i].Name)
		} else {
			exchangesTest = append(exchangesTest, configs[i].Name)
		}
	}

	for i := range exchanges {
		for j := range symbols {
			table[exchanges[i]+":"+symbols[j]] = make(map[string]domain.Exchange)
		}
	}

	for i := range exchangesTest {
		for j := range symbolsTest {
			table[exchangesTest[i]+":"+symbolsTest[j]] = make(map[string]domain.Exchange)
		}
	}

	return &storage{
		table: table,
		mu:    sync.Mutex{},
	}
}

func (st *storage) Write(exchange domain.Exchange) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.table[exchange.Exchange+":"+exchange.Symbol][exchange.ID] = exchange
}

func (st *storage) GetAll(exchanges, pairNames []string) []domain.Exchange {
	st.mu.Lock()
	defer st.mu.Unlock()

	var result []domain.Exchange

	for i := range exchanges {
		for j := range pairNames {
			for _, v := range st.table[exchanges[i]+":"+pairNames[j]] {
				result = append(result, v)
			}
		}
	}

	return result
}

func (st *storage) GetByPeriod(exchanges []string, symbol string, period int64) []domain.Exchange {
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.Lock()
	defer st.mu.Unlock()

	for i := range exchanges {
		delete(st.table[exchanges[i].Exchange+":"+exchanges[i].Symbol], exchanges[i].ID)
	}
}
