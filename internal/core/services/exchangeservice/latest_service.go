package exchangeservice

import (
	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetLatestSymbol(symbol string) (domain.Exchange, error) {
	var exchanges = exchangeServ.exchangeRepository.GetExchangesBySymbol(symbol)
	var result domain.Exchange

	results, err := exchangeServ.redisRepository.GetLatestSymbol(exchanges, symbol)
	if err != nil {
		return result, domain.ErrorNotFound
	}

	if len(results) == 0 {
		return result, domain.ErrorNotFound
	}

	result = GetLatest(results)

	return result, nil
}

func (exchangeServ *exchangeService) GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error) {
	result, err := exchangeServ.redisRepository.GetLatestExchangeSymbol(exchange, symbol)
	if err != nil {
		return result, domain.ErrorNotFound
	}

	return result, nil
}

func GetLatest(exchanges []domain.Exchange) domain.Exchange {
	var result domain.Exchange
	var latest int64

	for i := range exchanges {
		if i == 0 {
			latest = exchanges[i].Timestamp
			result = exchanges[i]
		} else if latest < exchanges[i].Timestamp {
			latest = exchanges[i].Timestamp
			result = exchanges[i]
		}
	}

	return result
}
