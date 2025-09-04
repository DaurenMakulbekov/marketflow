package exchangeservice

import (
	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetLatestPrice(symbol string) (domain.Exchange, error) {
	var result domain.Exchange
	var exchanges = exchangeServ.exchangeRepository.GetExchangesBySymbol(symbol)

	if len(exchanges) == 0 {
		return result, domain.ErrorBadRequest
	}

	results, err := exchangeServ.redisRepository.GetLatestPrice(exchanges, symbol)
	if err != nil {
		//return result, domain.ErrorNotFound
	}

	if len(results) == 0 {
		var resStorage = exchangeServ.storage.GetLatest(exchanges, symbol)

		if len(results) == 0 && len(resStorage) == 0 {
			return result, domain.ErrorNotFound
		}
		results = resStorage
	}

	result = GetLatest(results)

	return result, nil
}

func (exchangeServ *exchangeService) GetLatestExchangePrice(exchange, symbol string) (domain.Exchange, error) {
	var resExchange bool = exchangeServ.exchangeRepository.CheckExchange(exchange)
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resExchange == false || resSymbol == false {
		return domain.Exchange{}, domain.ErrorBadRequest
	}

	result, err := exchangeServ.redisRepository.GetLatestExchangePrice(exchange, symbol)
	if err != nil {
		//return result, domain.ErrorNotFound
	}

	if result.Price == 0 {
		var resStorage = exchangeServ.storage.GetLatestByExchange(exchange, symbol)

		if result.Price == 0 && len(resStorage) == 0 {
			return result, domain.ErrorNotFound
		}
		result = GetLatest(resStorage)
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
