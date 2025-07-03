package exchangeservice

import "marketflow/internal/core/ports"

type exchangeService struct {
	exchangeRepository ports.ExchangeRepository
}

func NewExchangeService(exchangeRepo ports.ExchangeRepository) *exchangeService {
	return &exchangeService{
		exchangeRepository: exchangeRepo,
	}
}

func (exchangeServ *exchangeService) GetData() {
	exchangeServ.exchangeRepository.GetData()
}
