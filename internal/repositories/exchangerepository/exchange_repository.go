package exchangerepository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"marketflow/internal/core/domain"
	"math/rand/v2"
	"net"
	"os"

	"time"
)

type exchangeRepository struct {
}

func NewExchangeRepository() *exchangeRepository {
	return &exchangeRepository{}
}

func (exchangeRepo *exchangeRepository) GetFromExchange(exchange string) <-chan string {
	var out = make(chan string)

	go func() {
		defer close(out)

		conn, err := net.Dial("tcp", exchange)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer conn.Close()

		var scanner = bufio.NewScanner(conn)

		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out
}

func (exchangeRepo *exchangeRepository) Generator() <-chan string {
	var out = make(chan string)
	var pairNames = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	go func() {
		defer close(out)

		for {
			for i := range pairNames {
				var exchange = domain.Exchange{
					Symbol:    pairNames[i],
					Timestamp: time.Now().Unix(),
				}

				if pairNames[i] == "BTCUSDT" {
					exchange.Price = (rand.Float64() * 6000) + 97000
				} else if pairNames[i] == "DOGEUSDT" {
					exchange.Price = (rand.Float64() * 0.07) + 0.28
				} else if pairNames[i] == "TONUSDT" {
					exchange.Price = (rand.Float64() * 1.1) + 3.4
				} else if pairNames[i] == "SOLUSDT" {
					exchange.Price = (rand.Float64() * 85) + 197
				} else if pairNames[i] == "ETHUSDT" {
					exchange.Price = (rand.Float64() * 470) + 2627
				}

				result, _ := json.Marshal(exchange)
				out <- string(result)
				time.Sleep(22 * time.Millisecond)
			}
		}
	}()

	return out
}
