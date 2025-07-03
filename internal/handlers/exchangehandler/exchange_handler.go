package exchangehandler

import (
	"marketflow/internal/core/ports"
	"net/http"
)

type exchangeHandler struct {
	exchangeService ports.ExchangeService
}

func NewExchangeHandler(exchangeServ ports.ExchangeService) *exchangeHandler {
	return &exchangeHandler{
		exchangeService: exchangeServ,
	}
}

func (exchangeHandl *exchangeHandler) LiveModeHandler(w http.ResponseWriter, req *http.Request) {
	exchangeHandl.exchangeService.GetData()
}
