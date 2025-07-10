package postgresrepository

import (
	"database/sql"
	"fmt"
	"marketflow/internal/core/domain"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(DB *sql.DB) *postgresRepository {
	return &postgresRepository{
		db: DB,
	}
}

func (postgresRepo *postgresRepository) Write(exchange domain.Exchanges) error {
	var query = `INSERT INTO exchanges (pair_name, exchange, timestamp, average_price, min_price, max_price) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := postgresRepo.db.Exec(query, exchange.PairName, exchange.Exchange, exchange.Timestamp, exchange.AvgPrice, exchange.MinPrice, exchange.MaxPrice)
	if err != nil {
		return fmt.Errorf("Error Write exchange data: %v", err)
	}

	return nil
}
