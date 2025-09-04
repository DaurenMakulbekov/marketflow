package exchangehandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"marketflow/internal/core/ports"
	"net/http"
	"os"
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
	exchangeHandl.exchangeService.LiveMode()
}

func (exchangeHandl *exchangeHandler) TestModeHandler(w http.ResponseWriter, req *http.Request) {
	exchangeHandl.exchangeService.TestMode()
}

func (exchangeHandl *exchangeHandler) SystemStatusHandler(w http.ResponseWriter, req *http.Request) {
	var result = exchangeHandl.exchangeService.GetSystemStatus()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(result)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err_)
		return
	}

	logger.Info("System Status", "method", "GET", "status", 200)
}
