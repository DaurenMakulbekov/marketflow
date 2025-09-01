package exchangeservice

import (
	"marketflow/internal/core/domain"
	"strconv"
	"strings"
	"time"
)

func (exchangeServ *exchangeService) GetHighestPrice(symbol string) (domain.PriceSymbol, error) {
	var result bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if result == false {
		return domain.PriceSymbol{}, domain.ErrorBadRequest
	}

	exchange, err := exchangeServ.postgresRepository.GetHighestPrice(symbol)
	if err != nil {
		return exchange, err
	}

	exchange.Symbol = symbol

	return exchange, nil
}

func (exchangeServ *exchangeService) GetHighestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error) {
	var resExchange bool = exchangeServ.exchangeRepository.CheckExchange(exchange)
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resExchange == false || resSymbol == false {
		return domain.PriceExchangeSymbol{}, domain.ErrorBadRequest
	}

	result, err := exchangeServ.postgresRepository.GetHighestExchangePrice(exchange, symbol)
	if err != nil {
		return result, err
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetHighestPriceByPeriod(symbol, period string) (domain.PriceSymbol, error) {
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resSymbol == false {
		return domain.PriceSymbol{}, domain.ErrorBadRequest
	}

	var result domain.PriceSymbol
	var err error

	s, err := time.ParseDuration(period)

	if err != nil || strings.Contains(period, "-") {
		return domain.PriceSymbol{}, domain.ErrorBadRequest
	}

	var timeNow = time.Now()
	var t = timeNow.Add(-time.Duration(s.Seconds()) * time.Second)
	var timestamp = t.UnixMilli()

	if timeNow.UnixMilli()-timestamp >= 60000 {
		result, err = exchangeServ.postgresRepository.GetHighestPriceByPeriod(symbol, t)
		if err != nil {
			return result, domain.ErrorNotFound
		}
	} else {
		var exchanges = exchangeServ.exchangeRepository.GetExchangesBySymbol(symbol)
		var id = strconv.FormatInt(timestamp, 10)

		res, err := exchangeServ.redisRepository.GetPriceByPeriod(exchanges, symbol, id)
		if err != nil {
			return result, domain.ErrorNotFound
		}

		var res1 = GetHighest(res)
		result.Price = res1.Price
	}

	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetHighestExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error) {
	var resExchange bool = exchangeServ.exchangeRepository.CheckExchange(exchange)
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resExchange == false || resSymbol == false {
		return domain.PriceExchangeSymbol{}, domain.ErrorBadRequest
	}

	var result domain.PriceExchangeSymbol
	var err error

	s, err := time.ParseDuration(period)

	if err != nil || strings.Contains(period, "-") {
		return domain.PriceExchangeSymbol{}, domain.ErrorBadRequest
	}

	var timeNow = time.Now()
	var t = timeNow.Add(-time.Duration(s.Seconds()) * time.Second)
	var timestamp = t.UnixMilli()

	if timeNow.UnixMilli()-timestamp >= 60000 {
		result, err = exchangeServ.postgresRepository.GetHighestExchangePriceByPeriod(exchange, symbol, t)
		if err != nil {
			return result, domain.ErrorNotFound
		}
	} else {
		var id = strconv.FormatInt(timestamp, 10)

		res, err := exchangeServ.redisRepository.GetExchangePriceByPeriod(exchange, symbol, id)
		if err != nil {
			return result, domain.ErrorNotFound
		}

		var res1 = GetHighest(res)
		result.Price = res1.Price
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func GetHighest(exchanges []domain.Exchange) domain.Exchange {
	var result domain.Exchange
	var max float64

	for i := range exchanges {
		if i == 0 {
			max = exchanges[i].Price
			result = exchanges[i]
		} else if max < exchanges[i].Price {
			max = exchanges[i].Price
			result = exchanges[i]
		}
	}

	return result
}
