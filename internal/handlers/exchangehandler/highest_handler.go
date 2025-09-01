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

func (exchangeHandl *exchangeHandler) HighestPriceHandler(w http.ResponseWriter, req *http.Request) {
	var symbol string = req.PathValue("symbol")
	var period string = req.URL.Query().Get("period")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var result domain.PriceSymbol
	var err error

	if len(period) > 0 {
		result, err = exchangeHandl.exchangeService.GetHighestPriceByPeriod(symbol, period)
	} else {
		result, err = exchangeHandl.exchangeService.GetHighestPrice(symbol)
	}

	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			PrintErrorMessage(w, req, http.StatusNotFound, "Not Found")
			logger.Error("Not Found", "method", "GET", "status", 404)
		} else {
			PrintErrorMessage(w, req, http.StatusBadRequest, "Incorrect input")
			logger.Error("Incorrect Input", "method", "GET", "status", 400)
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

	logger.Info("Retrieve a specific symbol", "method", "GET", "status", 200)
}

func (exchangeHandl *exchangeHandler) HighestExchangePriceHandler(w http.ResponseWriter, req *http.Request) {
	var exchange string = req.PathValue("exchange")
	var symbol string = req.PathValue("symbol")
	var period string = req.URL.Query().Get("period")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var result domain.PriceExchangeSymbol
	var err error

	if len(period) > 0 {
		result, err = exchangeHandl.exchangeService.GetHighestExchangePriceByPeriod(exchange, symbol, period)
	} else {
		result, err = exchangeHandl.exchangeService.GetHighestExchangePrice(exchange, symbol)
	}

	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			PrintErrorMessage(w, req, http.StatusNotFound, "Not Found")
			logger.Error("Not Found", "method", "GET", "status", 404)
		} else {
			PrintErrorMessage(w, req, http.StatusBadRequest, "Incorrect input")
			logger.Error("Incorrect Input", "method", "GET", "status", 400)
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

	logger.Info("Retrieve a specific symbol", "method", "GET", "status", 200)
}

func PrintErrorMessage(w http.ResponseWriter, req *http.Request, h int, s string) {
	w.WriteHeader(h)
	w.Header().Set("Content-Type", "application/json")
	m := map[string]string{"error": s}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}
