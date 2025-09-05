package exchangeservice

import (
	"strconv"
	"strings"
	"time"

	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetLowestPrice(symbol string) (domain.PriceSymbol, error) {
	var result bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if result == false {
		return domain.PriceSymbol{}, domain.ErrorBadRequest
	}

	exchange, err := exchangeServ.postgresRepository.GetLowestPrice(symbol)
	if err != nil {
		return exchange, err
	}

	exchange.Symbol = symbol

	return exchange, nil
}

func (exchangeServ *exchangeService) GetLowestExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error) {
	var resExchange bool = exchangeServ.exchangeRepository.CheckExchange(exchange)
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resExchange == false || resSymbol == false {
		return domain.PriceExchangeSymbol{}, domain.ErrorBadRequest
	}

	result, err := exchangeServ.postgresRepository.GetLowestExchangePrice(exchange, symbol)
	if err != nil {
		return result, err
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetLowestPriceByPeriod(symbol, period string) (domain.PriceSymbol, error) {
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

	timeNow := time.Now()
	t := timeNow.Add(-time.Duration(s.Seconds()) * time.Second)
	timestamp := t.UnixMilli()

	if timeNow.UnixMilli()-timestamp >= 60000 {
		result, err = exchangeServ.postgresRepository.GetLowestPriceByPeriod(symbol, t)
		if err != nil {
			return result, domain.ErrorNotFound
		}
	} else {
		exchanges := exchangeServ.exchangeRepository.GetExchangesBySymbol(symbol)
		id := strconv.FormatInt(timestamp, 10)

		resRedis, err := exchangeServ.redisRepository.GetPriceByPeriod(exchanges, symbol, id)
		if err != nil {
			// return result, domain.ErrorNotFound
		}

		resStorage := exchangeServ.storage.GetByPeriod(exchanges, symbol, timestamp)

		if len(resRedis) == 0 && len(resStorage) == 0 {
			return result, domain.ErrorNotFound
		}

		resRedis = append(resRedis, resStorage...)

		res1 := GetLowest(resRedis)
		result.Price = res1.Price
	}

	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetLowestExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error) {
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

	timeNow := time.Now()
	t := timeNow.Add(-time.Duration(s.Seconds()) * time.Second)
	timestamp := t.UnixMilli()

	if timeNow.UnixMilli()-timestamp >= 60000 {
		result, err = exchangeServ.postgresRepository.GetLowestExchangePriceByPeriod(exchange, symbol, t)
		if err != nil {
			return result, domain.ErrorNotFound
		}
	} else {
		id := strconv.FormatInt(timestamp, 10)

		resRedis, err := exchangeServ.redisRepository.GetExchangePriceByPeriod(exchange, symbol, id)
		if err != nil {
			// return result, domain.ErrorNotFound
		}

		resStorage := exchangeServ.storage.GetByExchangePeriod(exchange, symbol, timestamp)

		if len(resRedis) == 0 && len(resStorage) == 0 {
			return result, domain.ErrorNotFound
		}

		resRedis = append(resRedis, resStorage...)

		res1 := GetLowest(resRedis)
		result.Price = res1.Price
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func GetLowest(exchanges []domain.Exchange) domain.Exchange {
	var result domain.Exchange
	var min float64

	for i := range exchanges {
		if i == 0 {
			min = exchanges[i].Price
			result = exchanges[i]
		} else if min > exchanges[i].Price {
			min = exchanges[i].Price
			result = exchanges[i]
		}
	}

	return result
}
