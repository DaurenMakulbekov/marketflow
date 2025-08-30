package exchangeservice

import (
	"marketflow/internal/core/domain"
)

func (exchangeServ *exchangeService) GetHighestSymbol(symbol string) (domain.Price, error) {
	exchange, err := exchangeServ.postgresRepository.GetHighestSymbol(symbol)
	if err != nil {
		return exchange, err
	}

	return exchange, nil
}

/*  */
//var t = time.Now()

//var period string = "4m"

//s, _ := time.ParseDuration(period)

//var t1 = t.Add(-time.Duration(s.Minutes()) * time.Minute)

//result, err := GetHighestSymbol(db, t1)
//if err != nil {
//	fmt.Println(err)
//} else {
//	fmt.Println(result)
//}
/*  */
