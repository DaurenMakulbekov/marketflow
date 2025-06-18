package main

import (
	"fmt"
	"database/sql"
	"os"
	"net/http"
	"log"
	"context"
	"github.com/redis/go-redis/v9"
	_ "github.com/jackc/pgx/v5/stdlib"
)


func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

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
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Redis connection was refused")
	}

	fmt.Println(status)

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /", handler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
