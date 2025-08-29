package redisrepository

import (
	"marketflow/internal/core/domain"
	"strconv"
)

func (redisRepo *redisRepository) GetLatestSymbol(exchange, symbol string) (domain.Exchange, error) {
	res, err := redisRepo.rdb.XRevRangeN(redisRepo.ctx, exchange+":"+symbol, "+", "-", 1).Result()
	if err != nil {
		return domain.Exchange{}, err
	}

	price, _ := strconv.ParseFloat(res[1].Values["price"].(string), 64)
	timestamp, _ := strconv.ParseInt(res[1].Values["timestamp"].(string), 10, 64)

	var result = domain.Exchange{
		ID:        res[0].ID,
		Exchange:  res[1].Values["exchange"].(string),
		Symbol:    res[1].Values["symbol"].(string),
		Price:     price,
		Timestamp: timestamp,
	}

	return result, nil
}
