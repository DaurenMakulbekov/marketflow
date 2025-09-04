package exchangeservice

import "marketflow/internal/core/domain"

func (exchangeServ *exchangeService) GetSystemStatus() domain.SystemStatus {
	var result domain.SystemStatus

	var err = exchangeServ.redisRepository.CheckConnection()
	if err != nil {
		result.Redis = "Connection failed"
	}
	result.Redis = "Connected"

	var resultExchanges = exchangeServ.exchangeRepository.CheckConnection()

	result.Exchange1 = resultExchanges[0]
	result.Exchange2 = resultExchanges[1]
	result.Exchange3 = resultExchanges[2]

	return result
}
