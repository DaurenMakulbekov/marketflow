//package postgresrepository
//
//import (
//	"database/sql"
//	"fmt"
//	"marketflow/internal/core/domain"
//)
//
//func (postgresRepo *postgresRepository) GetHighestSymbol(symbol string) (domain.Exchanges, error) {
//	var exchange domain.Exchanges
//
//	row := postgresRepo.db.QueryRow("SELECT pair_name, max(max_price) FROM exchanges WHERE pair_name = $1 GROUP BY pair_name", symbol)
//	if err := row.Scan(&exchange.PairName, &exchange.MaxPrice); err != nil {
//		if err == sql.ErrNoRows {
//			return exchange, domain.ErrorNotFound //fmt.Errorf("Error: exchange by symbol %s: no such symbol", symbol)
//		}
//		return exchange, fmt.Errorf("Error: exchange by symbol %s: %v", symbol, err)
//	}
//
//	return exchange, nil
//}
