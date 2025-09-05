package main

import (
	"context"
	"errors"
	"fmt"
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

	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func helpMessage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  marketflow [--port <N>]")
	fmt.Println("  marketflow --help")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --port N  Port number")
	flag.PrintDefaults()
}

func main() {
	var port = flag.String("port", "8080", "")
	flag.Usage = helpMessage
	flag.Parse()

	number, _ := strconv.Atoi(*port)
	if number < 1024 || number > 49151 {
		fmt.Fprintln(os.Stderr, "Error: incorrect port number")
		os.Exit(1)
	}

	var config = config.NewAppConfig()
	var ctx = context.Background()

	var redisRepository = redisrepository.NewRedisRepository(config.Redis, ctx)
	var postgresRepository = postgresrepository.NewPostgresRepository(config.DB)
	var storageRepository = storage.NewStorage()
	var exchangeRepos = exchangerepository.NewExchangeRepository(config.Exchanges)
	var exchangeService = exchangeservice.NewExchangeService(exchangeRepos, redisRepository, postgresRepository, storageRepository)
	var exchangeHandler = exchangehandler.NewExchangeHandler(exchangeService)

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /prices/latest/{symbol}", exchangeHandler.LatestPriceHandler)
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", exchangeHandler.LatestExchangePriceHandler)
	mux.HandleFunc("GET /prices/highest/{symbol}", exchangeHandler.HighestPriceHandler)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", exchangeHandler.HighestExchangePriceHandler)
	mux.HandleFunc("GET /prices/lowest/{symbol}", exchangeHandler.LowestPriceHandler)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}", exchangeHandler.LowestExchangePriceHandler)
	mux.HandleFunc("GET /prices/average/{symbol}", exchangeHandler.AveragePriceHandler)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}", exchangeHandler.AverageExchangePriceHandler)
	mux.HandleFunc("POST /mode/live", exchangeHandler.LiveModeHandler)
	mux.HandleFunc("POST /mode/test", exchangeHandler.TestModeHandler)
	mux.HandleFunc("GET /health", exchangeHandler.SystemStatusHandler)

	var server = &http.Server{
		Addr:    ":" + *port,
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
