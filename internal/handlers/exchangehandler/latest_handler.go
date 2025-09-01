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

func (exchangeHandl *exchangeHandler) LatestPriceHandler(w http.ResponseWriter, req *http.Request) {
	var symbol string = req.PathValue("symbol")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	exchange, err := exchangeHandl.exchangeService.GetLatestPrice(symbol)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			PrintErrorMessage(w, req, http.StatusNotFound, "Not Found")
			logger.Error("Not Found", "method", "GET", "status", 404)
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

func (exchangeHandl *exchangeHandler) LatestExchangePriceHandler(w http.ResponseWriter, req *http.Request) {
	var exchange string = req.PathValue("exchange")
	var symbol string = req.PathValue("symbol")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	result, err := exchangeHandl.exchangeService.GetLatestExchangePrice(exchange, symbol)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			PrintErrorMessage(w, req, http.StatusNotFound, "Not Found")
			logger.Error("Not Found", "method", "GET", "status", 404)
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
