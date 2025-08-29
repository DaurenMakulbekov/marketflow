package exchangeservice

import (
	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetLatestSymbol(symbol string) (domain.Exchange, error) {
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}
	var result []domain.Exchange

	for i := range exchanges {
		exchange, err := exchangeServ.redisRepository.GetLatestSymbol(exchanges[i], symbol)
		if err != nil {
			return exchange, err
		}

		result = append(result, exchange)
	}

	var exchange = GetLatest(result)

	return exchange, nil
}

func (exchangeServ *exchangeService) GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error) {
	result, err := exchangeServ.redisRepository.GetLatestSymbol(exchange, symbol)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetLatest(exchanges []domain.Exchange) domain.Exchange {
	var result domain.Exchange
	var latest int64

	for i := range exchanges {
		if i == 0 {
			latest = exchanges[i].Timestamp
		} else if latest > exchanges[i].Timestamp {
			latest = exchanges[i].Timestamp
			result = exchanges[i]
		}
	}

	return result
}
