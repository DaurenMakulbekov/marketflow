package exchangerepository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"marketflow/internal/core/domain"
	"marketflow/internal/infrastructure/config"
	"math/rand/v2"
	"net"
	"slices"
	"sync"

	"time"
)

type exchangeRepository struct {
	table         map[string]*config.Exchange
	done          chan bool
	doneTest      chan bool
	doneWG        sync.WaitGroup
	exchanges     []string
	pairNames     []string
	exchangesTest []string
	pairNamesTest []string
}

func NewExchangeRepository(configs []*config.Exchange) *exchangeRepository {
	var table = make(map[string]*config.Exchange)
	var done = make(chan bool)
	var doneTest = make(chan bool)
	var exchanges []string
	var pairNames = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}
	var exchangesTest = []string{"exchange1_test", "exchange2_test", "exchange3_test"}
	var pairNamesTest = []string{"BTCUSDT_test", "DOGEUSDT_test", "TONUSDT_test", "SOLUSDT_test", "ETHUSDT_test"}

	for i := range configs {
		exchanges = append(exchanges, configs[i].Name)
	}

	for i := range exchanges {
		table[exchanges[i]] = configs[i]
	}

	return &exchangeRepository{
		table:         table,
		done:          done,
		doneTest:      doneTest,
		exchanges:     exchanges,
		pairNames:     pairNames,
		exchangesTest: exchangesTest,
		pairNamesTest: pairNamesTest,
	}
}

func (exchangeRepo *exchangeRepository) Close() {
	for i := 0; i < 3; i++ {
		exchangeRepo.done <- true
	}
	close(exchangeRepo.done)
}

func (exchangeRepo *exchangeRepository) CloseTest() {
	for i := 0; i < 3; i++ {
		exchangeRepo.doneTest <- true
	}
}

func (exchangeRepo *exchangeRepository) Stop(ctx context.Context) {
	fmt.Println("Waiting for exchange repository to finish")
	var done = make(chan struct{})

	go func() {
		exchangeRepo.doneWG.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context done earlier")
	case <-done:
		fmt.Println("Exchange repository finished")
	}
}

func Connect(exchange *config.Exchange) (net.Conn, error) {
	conn, err := net.Dial("tcp", exchange.Host+":"+exchange.Port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (exchangeRepo *exchangeRepository) ReConnect(exchange string, ticker *time.Ticker, connect chan<- net.Conn, done chan<- bool) {
	go func() {
		for {
			select {
			case <-exchangeRepo.done:
				done <- true
				return
			case <-ticker.C:
				conn, err := Connect(exchangeRepo.table[exchange])
				if err != nil {
					//fmt.Fprintf(os.Stderr, "Failed to reconnect to %s\n", exchange)
				} else {
					log.Printf("Connected to %s\n", exchange)
					connect <- conn
					return
				}
			}
		}
	}()
}

func (exchangeRepo *exchangeRepository) GetFromExchange(exchange string) <-chan string {
	var out = make(chan string)

	go func() {
		defer close(out)

		var connect = make(chan net.Conn)
		defer close(connect)
		var done = make(chan bool)
		defer close(done)

		for {
			conn, err := Connect(exchangeRepo.table[exchange])
			if err != nil {
				log.Printf("Failed to connect to %s\n", exchange)

				var ticker = time.NewTicker(time.Second)
				defer ticker.Stop()

				exchangeRepo.ReConnect(exchange, ticker, connect, done)

				select {
				case conn1 := <-connect:
					conn = conn1
				case <-done:
					return
				}
			}
			var scanner = bufio.NewScanner(conn)

			for scanner.Scan() {
				select {
				case <-exchangeRepo.done:
					return
				default:
					out <- scanner.Text()
				}
			}
			conn.Close()
		}
	}()

	return out
}

func (exchangeRepo *exchangeRepository) Generator() <-chan string {
	var out = make(chan string)
	var pairNames = exchangeRepo.pairNamesTest

	go func() {
		defer close(out)

		for {
			select {
			case <-exchangeRepo.doneTest:
				return
			default:
				for i := range pairNames {
					var exchange = domain.Exchange{
						Symbol:    pairNames[i],
						Timestamp: time.Now().UnixMilli(),
					}

					if pairNames[i] == "BTCUSDT_test" {
						exchange.Price = (rand.Float64() * 6000) + 97000
					} else if pairNames[i] == "DOGEUSDT_test" {
						exchange.Price = (rand.Float64() * 0.07) + 0.28
					} else if pairNames[i] == "TONUSDT_test" {
						exchange.Price = (rand.Float64() * 1.1) + 3.4
					} else if pairNames[i] == "SOLUSDT_test" {
						exchange.Price = (rand.Float64() * 85) + 197
					} else if pairNames[i] == "ETHUSDT_test" {
						exchange.Price = (rand.Float64() * 470) + 2627
					}

					result, _ := json.Marshal(exchange)
					out <- string(result)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return out
}

func (exchangeRepo *exchangeRepository) GetExchanges() ([]string, []string) {
	return exchangeRepo.exchanges, exchangeRepo.pairNames
}

func (exchangeRepo *exchangeRepository) GetExchangesTest() ([]string, []string) {
	return exchangeRepo.exchangesTest, exchangeRepo.pairNamesTest
}

func (exchangeRepo *exchangeRepository) GetExchangesBySymbol(symbol string) []string {
	var result []string

	var res = slices.Contains(exchangeRepo.pairNames, symbol)
	if res == true {
		result = exchangeRepo.exchanges
		return result
	}

	var resTest = slices.Contains(exchangeRepo.pairNamesTest, symbol)
	if resTest == true {
		result = exchangeRepo.exchangesTest
		return result
	}

	return result
}

func (exchangeRepo *exchangeRepository) CheckSymbol(symbol string) bool {
	var res = slices.Contains(exchangeRepo.pairNames, symbol)
	if res == true {
		return true
	}

	var resTest = slices.Contains(exchangeRepo.pairNamesTest, symbol)
	if resTest == true {
		return true
	}

	return false
}

func (exchangeRepo *exchangeRepository) CheckExchange(exchange string) bool {
	var res = slices.Contains(exchangeRepo.exchanges, exchange)
	if res == true {
		return true
	}

	var resTest = slices.Contains(exchangeRepo.exchangesTest, exchange)
	if resTest == true {
		return true
	}

	return false
}

func (exchangeRepo *exchangeRepository) CheckConnection() []string {
	var result []string

	for i := range exchangeRepo.exchanges {
		var config = exchangeRepo.table[exchangeRepo.exchanges[i]]

		conn, err := Connect(config)
		if err != nil {
			result = append(result, "Connection failed")
		} else {
			result = append(result, "Connected")
			conn.Close()
		}
	}

	return result
}
