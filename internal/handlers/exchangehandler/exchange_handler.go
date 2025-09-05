package exchangehandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"marketflow/internal/core/ports"
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

	m := map[string]string{"Live mode": "active"}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(m)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err_)
		return
	}

	logger.Info("Live mode", "method", "POST", "status", 200)
}

func (exchangeHandl *exchangeHandler) TestModeHandler(w http.ResponseWriter, req *http.Request) {
	exchangeHandl.exchangeService.TestMode()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	m := map[string]string{"Test mode": "active"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(m)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err_)
		return
	}

	logger.Info("Test mode", "method", "POST", "status", 200)
}

func (exchangeHandl *exchangeHandler) SystemStatusHandler(w http.ResponseWriter, req *http.Request) {
	result := exchangeHandl.exchangeService.GetSystemStatus()

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
