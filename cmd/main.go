package main

import (
	"context"
	"errors"
	"time"

	"log"
	"marketflow/internal/core/services/exchangeservice"
	"marketflow/internal/handlers/exchangehandler"
	"marketflow/internal/infrastructure/config"
	"marketflow/internal/repositories/exchangerepository"
	"marketflow/internal/repositories/postgresrepository"
	"marketflow/internal/repositories/redisrepository"
	"marketflow/internal/repositories/storage"
	"net/http"

	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var config = config.NewAppConfig()

	var ctx = context.Background()

	var redisRepository = redisrepository.NewRedisRepository(config.Redis, ctx)
	var postgresRepository = postgresrepository.NewPostgresRepository(config.DB)
	var storageRepository = storage.NewStorage()
	var exchangeRepos = exchangerepository.NewExchangeRepository(config.Exchanges)
	var exchangeService = exchangeservice.NewExchangeService(exchangeRepos, redisRepository, postgresRepository, storageRepository)
	var exchangeHandler = exchangehandler.NewExchangeHandler(exchangeService)

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /prices/latest/{symbol}", exchangeHandler.LatestSymbolHandler)
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", exchangeHandler.LatestExchangeSymbolHandler)

	mux.HandleFunc("GET /prices/highest/{symbol}", exchangeHandler.HighestPriceHandler)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", exchangeHandler.HighestExchangePriceHandler)

	mux.HandleFunc("POST /mode/live", exchangeHandler.LiveModeHandler)
	mux.HandleFunc("POST /mode/test", exchangeHandler.TestModeHandler)

	var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signalCtxStop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen and serve returned error: %v", err)
		}
	}()

	<-signalCtx.Done()

	log.Println("Shutting down server...")
	exchangeRepos.Close()
	time.Sleep(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
	}

	log.Println("Server shutdown complete")
}
