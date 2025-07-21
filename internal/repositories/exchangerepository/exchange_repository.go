package exchangerepository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"marketflow/internal/core/domain"
	"marketflow/internal/core/ports"
	"math/rand/v2"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type exchangeRepository struct {
	redisRepository    ports.RedisRepository
	postgresRepository ports.PostgresRepository
}

func NewExchangeRepository(redisRepo ports.RedisRepository, postgresRepo ports.PostgresRepository) *exchangeRepository {
	return &exchangeRepository{
		redisRepository:    redisRepo,
		postgresRepository: postgresRepo,
	}
}

type Exchange struct {
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type ExchangeTest struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func GetFromExchange(exchange string) <-chan string {
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

func Distributor(exchangesAddress []string) []<-chan string {
	var out = make([]<-chan string, 3)

	for i := range exchangesAddress {
		out[i] = GetFromExchange(exchangesAddress[i])
	}

	return out
}

func Worker(in <-chan string, exchangeName string) <-chan string {
	var out = make(chan string)

	go func() {
		defer close(out)

		for i := range in {
			var result Exchange

			var decoder = json.NewDecoder(strings.NewReader(i))

			var err = decoder.Decode(&result)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error decode:", err)
			}

			var result1 = Exchange{
				Exchange:  exchangeName,
				Symbol:    result.Symbol,
				Price:     result.Price,
				Timestamp: result.Timestamp,
			}

			result2, _ := json.Marshal(result1)

			out <- string(result2)
		}
	}()

	return out
}

func Merger(ins ...<-chan string) <-chan string {
	var out = make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, in := range ins {
		go func(in <-chan string) {
			defer wg.Done()

			for i := range in {
				out <- i
			}
		}(in)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func CreateHashTable(exchanges, pairName []string) map[string]map[string]map[string]float64 {
	var m = make(map[string]map[string]map[string]float64, 3)

	for i := range exchanges {
		m[exchanges[i]] = make(map[string]map[string]float64, 5)
		for j := range pairName {
			m[exchanges[i]][pairName[j]] = make(map[string]float64, 4)
			m[exchanges[i]][pairName[j]]["min"] = 0
		}
	}

	return m
}

func (exchangeRepo *exchangeRepository) Aggregate(exchanges, pairName []string, m map[string]map[string]map[string]float64) {
	var length = exchangeRepo.redisRepository.LLen()

	for i := 0; i < length; i++ {
		var data Exchange
		var exchange = exchangeRepo.redisRepository.Read()
		var decode = json.NewDecoder(strings.NewReader(exchange))

		var err = decode.Decode(&data)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decode:", err)
		}

		m[data.Exchange][data.Symbol]["avg"] += data.Price

		if m[data.Exchange][data.Symbol]["min"] > data.Price || m[data.Exchange][data.Symbol]["min"] == 0 {
			m[data.Exchange][data.Symbol]["min"] = data.Price
		}
		if m[data.Exchange][data.Symbol]["max"] < data.Price {
			m[data.Exchange][data.Symbol]["max"] = data.Price
		}

		m[data.Exchange][data.Symbol]["count"]++
	}
}

func (exchangeRepo *exchangeRepository) Write(exchanges, pairName []string, m map[string]map[string]map[string]float64) {
	for i := range exchanges {
		for j := range pairName {
			var exchange = domain.Exchanges{
				PairName:  pairName[j],
				Exchange:  exchanges[i],
				Timestamp: time.Now(),
				AvgPrice:  m[exchanges[i]][pairName[j]]["avg"] / m[exchanges[i]][pairName[j]]["count"],
				MinPrice:  m[exchanges[i]][pairName[j]]["min"],
				MaxPrice:  m[exchanges[i]][pairName[j]]["max"],
			}

			var err = exchangeRepo.postgresRepository.Write(exchange)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (exchangeRepo *exchangeRepository) WriteToStorage(exchanges []string, ticker *time.Ticker, done chan bool) {
	var pairName = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var m = CreateHashTable(exchanges, pairName)
				exchangeRepo.Aggregate(exchanges, pairName, m)
				exchangeRepo.Write(exchanges, pairName, m)
			}
		}
	}()
}

func (exchangeRepo *exchangeRepository) LiveMode() {
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}
	var exchangesAddress = []string{"172.22.0.5:40101", "172.22.0.6:40102", "172.22.0.7:40103"}
	var outSlice = make([]<-chan string, 15)
	var index int = 0

	var out = Distributor(exchangesAddress)

	for i := range out {
		for j := 0; j < 5; j++ {
			outSlice[index] = Worker(out[i], exchanges[i])
			index++
		}
	}

	var merged = Merger(outSlice...)
	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)

	exchangeRepo.WriteToStorage(exchanges, ticker, done)

	for i := range merged {
		exchangeRepo.redisRepository.Write(i)
	}

	done <- true
	ticker.Stop()
}

func Generator() <-chan string {
	var out = make(chan string)
	var pairName = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	go func() {
		defer close(out)

		for {
			for i := range pairName {
				var exchange = Exchange{
					Symbol:    pairName[i],
					Timestamp: time.Now().Unix(),
				}

				if pairName[i] == "BTCUSDT" {
					exchange.Price = (rand.Float64() * 6000) + 97000
				} else if pairName[i] == "DOGEUSDT" {
					exchange.Price = (rand.Float64() * 0.07) + 0.28
				} else if pairName[i] == "TONUSDT" {
					exchange.Price = (rand.Float64() * 1.1) + 3.4
				} else if pairName[i] == "SOLUSDT" {
					exchange.Price = (rand.Float64() * 85) + 197
				} else if pairName[i] == "ETHUSDT" {
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

func (exchangeRepo *exchangeRepository) TestMode() {
	var exchanges = []string{"exchange_test"}
	var outSlice = make([]<-chan string, 10)
	var out = Generator()

	for i := 0; i < 5; i++ {
		outSlice[i] = Worker(out, "exchange_test")
	}

	var merged = Merger(outSlice...)
	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)

	exchangeRepo.WriteToStorage(exchanges, ticker, done)

	for i := range merged {
		exchangeRepo.redisRepository.Write(i)
	}

	done <- true
	ticker.Stop()
}
