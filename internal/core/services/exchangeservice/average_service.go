package exchangeservice

import (
	"strconv"
	"strings"
	"time"

	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetAveragePrice(symbol string) (domain.PriceSymbol, error) {
	var result bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if result == false {
		return domain.PriceSymbol{}, domain.ErrorBadRequest
	}

	exchange, err := exchangeServ.postgresRepository.GetAveragePrice(symbol)
	if err != nil {
		return exchange, err
	}

	exchange.Symbol = symbol

	return exchange, nil
}

func (exchangeServ *exchangeService) GetAverageExchangePrice(exchange, symbol string) (domain.PriceExchangeSymbol, error) {
	var resExchange bool = exchangeServ.exchangeRepository.CheckExchange(exchange)
	var resSymbol bool = exchangeServ.exchangeRepository.CheckSymbol(symbol)

	if resExchange == false || resSymbol == false {
		return domain.PriceExchangeSymbol{}, domain.ErrorBadRequest
	}

	result, err := exchangeServ.postgresRepository.GetAverageExchangePrice(exchange, symbol)
	if err != nil {
		return result, err
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetAverageExchangePriceByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error) {
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
		result, err = exchangeServ.postgresRepository.GetAverageExchangePriceByPeriod(exchange, symbol, t)
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

		result.Price = GetAverage(resRedis)
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func GetAverage(exchanges []domain.Exchange) float64 {
	var sum float64
	var count int = 0

	for i := range exchanges {
		sum += exchanges[i].Price
		count++
	}

	sum = sum / float64(count)

	return sum
}
