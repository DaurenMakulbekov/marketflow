package redisrepository

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"marketflow/internal/core/domain"
	"strconv"
)

func (redisRepo *redisRepository) GetLatestSymbol(exchanges []string, symbol string) ([]domain.Exchange, error) {
	var exchangesData []domain.Exchange
	var tx = redisRepo.rdb.TxPipeline()

	for i := range exchanges {
		tx.XRevRangeN(redisRepo.ctx, exchanges[i]+":"+symbol, "+", "-", 1)
	}

	cmds, err := tx.Exec(redisRepo.ctx)
	if err != nil {
		return []domain.Exchange{}, fmt.Errorf("Key not found in Redis cache")
	}

	for _, c := range cmds {
		var result = c.(*redis.XMessageSliceCmd).Val()

		for i := range result {
			price, _ := strconv.ParseFloat(result[i].Values["price"].(string), 64)
			timestamp, _ := strconv.ParseInt(result[i].Values["timestamp"].(string), 10, 64)

			var exchange = domain.Exchange{
				ID:        result[i].ID,
				Exchange:  result[i].Values["exchange"].(string),
				Symbol:    result[i].Values["symbol"].(string),
				Price:     price,
				Timestamp: timestamp,
			}

			exchangesData = append(exchangesData, exchange)
		}
	}

	return exchangesData, nil
}

func (redisRepo *redisRepository) GetLatestExchangeSymbol(exchange, symbol string) (domain.Exchange, error) {
	res, err := redisRepo.rdb.XRevRangeN(redisRepo.ctx, exchange+":"+symbol, "+", "-", 1).Result()
	if err != nil {
		return domain.Exchange{}, err
	}

	if len(res) == 0 {
		return domain.Exchange{}, domain.ErrorNotFound
	}

	price, _ := strconv.ParseFloat(res[0].Values["price"].(string), 64)
	timestamp, _ := strconv.ParseInt(res[0].Values["timestamp"].(string), 10, 64)

	var result = domain.Exchange{
		ID:        res[0].ID,
		Exchange:  res[0].Values["exchange"].(string),
		Symbol:    res[0].Values["symbol"].(string),
		Price:     price,
		Timestamp: timestamp,
	}

	return result, nil
}

func (redisRepo *redisRepository) GetPriceByPeriod(exchanges []string, symbol, period string) ([]domain.Exchange, error) {
	var exchangesData []domain.Exchange
	var tx = redisRepo.rdb.TxPipeline()

	for i := range exchanges {
		tx.XRange(redisRepo.ctx, exchanges[i]+":"+symbol, period+"-0", "+")
	}

	cmds, err := tx.Exec(redisRepo.ctx)
	if err != nil {
		return []domain.Exchange{}, fmt.Errorf("Key not found in Redis cache")
	}

	for _, c := range cmds {
		var result = c.(*redis.XMessageSliceCmd).Val()

		for i := range result {
			price, _ := strconv.ParseFloat(result[i].Values["price"].(string), 64)
			timestamp, _ := strconv.ParseInt(result[i].Values["timestamp"].(string), 10, 64)

			var exchange = domain.Exchange{
				ID:        result[i].ID,
				Exchange:  result[i].Values["exchange"].(string),
				Symbol:    result[i].Values["symbol"].(string),
				Price:     price,
				Timestamp: timestamp,
			}

			exchangesData = append(exchangesData, exchange)
		}
	}

	return exchangesData, nil
}

func (redisRepo *redisRepository) GetExchangePriceByPeriod(exchange, symbol, period string) ([]domain.Exchange, error) {
	var exchangesData []domain.Exchange
	var tx = redisRepo.rdb.TxPipeline()

	tx.XRange(redisRepo.ctx, exchange+":"+symbol, period+"-0", "+")

	cmds, err := tx.Exec(redisRepo.ctx)
	if err != nil {
		return []domain.Exchange{}, fmt.Errorf("Key not found in Redis cache")
	}

	for _, c := range cmds {
		var result = c.(*redis.XMessageSliceCmd).Val()

		for i := range result {
			price, _ := strconv.ParseFloat(result[i].Values["price"].(string), 64)
			timestamp, _ := strconv.ParseInt(result[i].Values["timestamp"].(string), 10, 64)

			var exchange = domain.Exchange{
				ID:        result[i].ID,
				Exchange:  result[i].Values["exchange"].(string),
				Symbol:    result[i].Values["symbol"].(string),
				Price:     price,
				Timestamp: timestamp,
			}

			exchangesData = append(exchangesData, exchange)
		}
	}

	return exchangesData, nil
}
