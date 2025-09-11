package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"marketflow/internal/core/services/exchangeservice"
	"marketflow/internal/handlers/exchangehandler"
	"marketflow/internal/infrastructure/config"
	"marketflow/internal/repositories/exchangerepository"
	"marketflow/internal/repositories/postgresrepository"
	"marketflow/internal/repositories/redisrepository"
	"marketflow/internal/repositories/storage"

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
	port := flag.String("port", "8080", "")
	flag.Usage = helpMessage
	flag.Parse()

	number, _ := strconv.Atoi(*port)
	if number < 1024 || number > 49151 {
		fmt.Fprintln(os.Stderr, "Error: incorrect port number")
		os.Exit(1)
	}

	config := config.NewAppConfig()
	ctx := context.Background()

	redisRepository := redisrepository.NewRedisRepository(config.Redis, ctx)
	postgresRepository := postgresrepository.NewPostgresRepository(config.DB)
	storageRepository := storage.NewStorage(config.Exchanges)
	exchangeRepos := exchangerepository.NewExchangeRepository(config.Exchanges)
	exchangeService := exchangeservice.NewExchangeService(exchangeRepos, redisRepository, postgresRepository, storageRepository)
	exchangeHandler := exchangehandler.NewExchangeHandler(exchangeService)
	exchangeService.LiveMode()

	mux := http.NewServeMux()

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

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	defer signalCtxStop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen and serve returned error: %v", err)
		}
	}()

	<-signalCtx.Done()

	log.Println("Shutting down server...")
	exchangeService.Close()
	time.Sleep(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
	}

	log.Println("Server shutdown complete")
}
