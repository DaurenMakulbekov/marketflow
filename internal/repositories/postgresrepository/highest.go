package postgresrepository

import (
	"database/sql"
	"fmt"
	"marketflow/internal/core/domain"
)

func (postgresRepo *postgresRepository) GetHighestSymbol(symbol string) (domain.Price, error) {
	var exchange domain.Price

	row := postgresRepo.db.QueryRow("SELECT max(max_price) FROM exchanges WHERE pair_name = $1", symbol)
	if err := row.Scan(&exchange.Price); err != nil {
		if err == sql.ErrNoRows {
			return exchange, domain.ErrorNotFound
		}
		return exchange, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
	}

	return exchange, nil
}
