package exchangeservice

import (
	"marketflow/internal/core/domain"
	"strconv"
	"strings"
	"time"
)

func (exchangeServ *exchangeService) GetHighestSymbol(symbol string) (domain.PriceSymbol, error) {
	exchange, err := exchangeServ.postgresRepository.GetHighestSymbol(symbol)
	if err != nil {
		return exchange, err
	}

	exchange.Symbol = symbol

	return exchange, nil
}

func (exchangeServ *exchangeService) GetHighestExchangeSymbol(exchange, symbol string) (domain.PriceExchangeSymbol, error) {
	result, err := exchangeServ.postgresRepository.GetHighestExchangeSymbol(exchange, symbol)
	if err != nil {
		return result, err
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetHighestSymbolByPeriod(symbol, period string) (domain.PriceSymbol, error) {
	var result domain.PriceSymbol
	var err error

	s, _ := time.ParseDuration(period)

	var timestamp time.Time

	if strings.Contains(period, "m") {
		timestamp = time.Now().Add(-time.Duration(s.Minutes()) * time.Minute)

		result, err = exchangeServ.postgresRepository.GetHighestSymbolByPeriod(symbol, timestamp)
		if err != nil {
			return result, domain.ErrorNotFound
		}
	} else if strings.Contains(period, "s") {
		timestamp = time.Now().Add(-time.Duration(s.Seconds()) * time.Second)

		var exchanges = exchangeServ.exchangeRepository.GetExchangesBySymbol(symbol)
		var id = strconv.FormatInt(timestamp.UnixMilli(), 10)

		res, err := exchangeServ.redisRepository.GetPriceByPeriod(exchanges, symbol, id)
		if err != nil {
			return result, err
		}

		var res1 = GetLatest(res)
		result.Price = res1.Price
	}

	result.Symbol = symbol

	return result, nil
}

func (exchangeServ *exchangeService) GetHighestExchangeSymbolByPeriod(exchange, symbol, period string) (domain.PriceExchangeSymbol, error) {
	s, _ := time.ParseDuration(period)

	var timestamp time.Time

	if strings.Contains(period, "m") {
		timestamp = time.Now().Add(-time.Duration(s.Minutes()) * time.Minute)
	} else if strings.Contains(period, "s") {
		timestamp = time.Now().Add(-time.Duration(s.Seconds()) * time.Second)
	}

	result, err := exchangeServ.postgresRepository.GetHighestExchangeSymbolByPeriod(exchange, symbol, timestamp)
	if err != nil {
		return result, domain.ErrorNotFound
	}

	result.Exchange = exchange
	result.Symbol = symbol

	return result, nil
}
