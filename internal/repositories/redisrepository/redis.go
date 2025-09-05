package redisrepository

import (
	"context"
	"fmt"
	"strconv"

	"marketflow/internal/core/domain"
	"marketflow/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	rdb    *redis.Client
	ctx    context.Context
	config *config.Redis
}

func NewRedisRepository(config *config.Redis, ctx context.Context) *redisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       0,
	})

	newRedisRepository := redisRepository{
		rdb:    rdb,
		ctx:    ctx,
		config: config,
	}

	return &newRedisRepository
}

func (redisRepo *redisRepository) Write(exchange domain.Exchange) error {
	_, err := redisRepo.rdb.XAdd(redisRepo.ctx, &redis.XAddArgs{
		Stream: exchange.Exchange + ":" + exchange.Symbol,
		Values: map[string]interface{}{
			"exchange":  exchange.Exchange,
			"symbol":    exchange.Symbol,
			"price":     exchange.Price,
			"timestamp": exchange.Timestamp,
		},
		ID: exchange.ID,
	}).Result()
	if err != nil {
		return fmt.Errorf("Failed to add key-value pair")
	}

	return nil
}

func (redisRepo *redisRepository) ReadAll(exchanges, pairNames []string) ([]domain.Exchange, error) {
	var exchangesData []domain.Exchange
	tx := redisRepo.rdb.TxPipeline()

	for i := range exchanges {
		for j := range pairNames {
			tx.XRange(redisRepo.ctx, exchanges[i]+":"+pairNames[j], "-", "+")
		}
	}

	cmds, err := tx.Exec(redisRepo.ctx)
	if err != nil {
		return []domain.Exchange{}, fmt.Errorf("Key not found in Redis cache")
	}

	for _, c := range cmds {
		result := c.(*redis.XMessageSliceCmd).Val()

		for i := range result {
			price, _ := strconv.ParseFloat(result[i].Values["price"].(string), 64)
			timestamp, _ := strconv.ParseInt(result[i].Values["timestamp"].(string), 10, 64)

			exchange := domain.Exchange{
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

func (redisRepo *redisRepository) DeleteAll(exchanges []domain.Exchange) error {
	tx := redisRepo.rdb.TxPipeline()

	for i := range exchanges {
		tx.XDel(redisRepo.ctx, exchanges[i].Exchange+":"+exchanges[i].Symbol, exchanges[i].ID)
	}

	_, err := tx.Exec(redisRepo.ctx)
	if err != nil {
		return fmt.Errorf("Error Delete: %v", err)
	}

	return nil
}

func (redisRepo *redisRepository) CheckConnection() error {
	_, err := redisRepo.rdb.Ping(redisRepo.ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisRepo *redisRepository) Reconnect() {
	newRedisRepository := NewRedisRepository(redisRepo.config, redisRepo.ctx)

	redisRepo.rdb = newRedisRepository.rdb
}

func (redisRepo *redisRepository) Close() {
	redisRepo.rdb.Close()

	fmt.Println("Redis Connection Closed")
}
