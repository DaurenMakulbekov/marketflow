package postgresrepository

import (
	"database/sql"
	"fmt"
	"marketflow/internal/core/domain"
	"time"
)

func (postgresRepo *postgresRepository) GetHighestSymbol(symbol string) (domain.PriceSymbol, error) {
	var exchange domain.PriceSymbol

	row := postgresRepo.db.QueryRow("SELECT max(max_price) FROM exchanges WHERE pair_name = $1", symbol)
	if err := row.Scan(&exchange.Price); err != nil {
		if err == sql.ErrNoRows {
			return exchange, domain.ErrorNotFound
		}
		return exchange, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
	}

	return exchange, nil
}

func (postgresRepo *postgresRepository) GetHighestExchangeSymbol(exchange, symbol string) (domain.PriceExchangeSymbol, error) {
	var result domain.PriceExchangeSymbol

	row := postgresRepo.db.QueryRow("SELECT max(max_price) FROM exchanges WHERE pair_name = $1 AND exchange = $2", symbol, exchange)
	if err := row.Scan(&result.Price); err != nil {
		if err == sql.ErrNoRows {
			return result, domain.ErrorNotFound
		}
		return result, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
	}

	return result, nil
}

func (postgresRepo *postgresRepository) GetHighestSymbolByPeriod(symbol string, period time.Time) (domain.PriceSymbol, error) {
	var exchange domain.PriceSymbol

	row := postgresRepo.db.QueryRow("SELECT max(max_price) FROM exchanges WHERE pair_name = $1 AND timestamp >= $2", symbol, period)
	if err := row.Scan(&exchange.Price); err != nil {
		if err == sql.ErrNoRows {
			return exchange, domain.ErrorNotFound
		}
		return exchange, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
	}

	return exchange, nil
}

func (postgresRepo *postgresRepository) GetHighestExchangeSymbolByPeriod(exchange, symbol string, period time.Time) (domain.PriceExchangeSymbol, error) {
	var result domain.PriceExchangeSymbol

	row := postgresRepo.db.QueryRow("SELECT max(max_price) FROM exchanges WHERE pair_name = $1 AND exchange = $2 AND timestamp >= $3", symbol, exchange, period)
	if err := row.Scan(&result.Price); err != nil {
		if err == sql.ErrNoRows {
			return result, domain.ErrorNotFound
		}
		return result, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
	}

	return result, nil
}
