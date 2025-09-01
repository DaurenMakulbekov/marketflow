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
		return result, domain.ErrorNotFound
	}

	if len(results) == 0 {
		return result, domain.ErrorNotFound
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
