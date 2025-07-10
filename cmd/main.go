package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"marketflow/internal/core/services/exchangeservice"
	"marketflow/internal/handlers/exchangehandler"
	"marketflow/internal/repositories/exchangerepository"
	"marketflow/internal/repositories/postgresrepository"
	"marketflow/internal/repositories/redisrepository"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func main() {
	var db *sql.DB
	var err error

	db, err = sql.Open("pgx", "user=user password=user host=db port=5432 database=marketflow sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	var ctx = context.Background()

	var rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Redis connection was refused")
	}

	fmt.Println(status)

	var redisRepository = redisrepository.NewRedisRepository(rdb, ctx)
	var postgresRepository = postgresrepository.NewPostgresRepository((db))
	var exchangeRepository = exchangerepository.NewExchangeRepository(redisRepository, postgresRepository)
	var exchangeService = exchangeservice.NewExchangeService(exchangeRepository)
	var exchangeHandler = exchangehandler.NewExchangeHandler(exchangeService)

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /mode/live", exchangeHandler.LiveModeHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
