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

func (postgresRepo *postgresRepository) Write(exchanges []domain.Exchanges) error {
	tx, err := postgresRepo.db.Begin()
	if err != nil {
		return fmt.Errorf("Error Transaction Begin: %v", err)
	}
	defer tx.Rollback()

	for i := range exchanges {
		var query = `INSERT INTO exchanges (pair_name, exchange, timestamp, average_price, min_price, max_price) VALUES($1, $2, $3, $4, $5, $6)`

		_, err := tx.Exec(query, exchanges[i].PairName, exchanges[i].Exchange, exchanges[i].Timestamp, exchanges[i].AvgPrice, exchanges[i].MinPrice, exchanges[i].MaxPrice)
		if err != nil {
			return fmt.Errorf("Error Write exchange data: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Error Transaction Commit: %v", err)
	}

	return nil
}
