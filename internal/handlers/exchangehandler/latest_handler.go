package exchangehandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"marketflow/internal/core/domain"
	"net/http"
	"os"
)

func (exchangeHandl *exchangeHandler) LatestSymbolHandler(w http.ResponseWriter, req *http.Request) {
	var symbol string = req.PathValue("symbol")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	exchange, err := exchangeHandl.exchangeService.GetLatestSymbol(symbol)
	if err != nil {
		if errors.Is(err, domain.ErrorBadRequest) {
			PrintErrorMessage(w, req, http.StatusBadRequest, "Incorrect input")
			logger.Error("Incorrect input", "method", "GET", "status", 400)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(exchange)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err_)
		return
	}

	logger.Info("Retrieve a specific symbol price", "method", "GET", "status", 200)
}

func (exchangeHandl *exchangeHandler) LatestExchangeSymbolHandler(w http.ResponseWriter, req *http.Request) {
	var exchange string = req.PathValue("exchange")
	var symbol string = req.PathValue("symbol")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	result, err := exchangeHandl.exchangeService.GetLatestExchangeSymbol(exchange, symbol)
	if err != nil {
		if errors.Is(err, domain.ErrorBadRequest) {
			PrintErrorMessage(w, req, http.StatusBadRequest, "Incorrect input")
			logger.Error("Incorrect input", "method", "GET", "status", 400)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(result)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err_)
		return
	}

	logger.Info("Retrieve a specific exchange symbol price", "method", "GET", "status", 200)
}
