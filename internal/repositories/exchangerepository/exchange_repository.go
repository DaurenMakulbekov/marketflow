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

func Worker(in <-chan string, exchangeName string) <-chan domain.Exchange {
	var out = make(chan domain.Exchange)

	go func() {
		defer close(out)

		for i := range in {
			var result domain.Exchange

			var decoder = json.NewDecoder(strings.NewReader(i))

			var err = decoder.Decode(&result)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error decode:", err)
			}

			var result1 = domain.Exchange{
				Exchange:  exchangeName,
				Symbol:    result.Symbol,
				Price:     result.Price,
				Timestamp: result.Timestamp,
			}

			out <- result1
		}
	}()

	return out
}

func Merger(ins ...<-chan domain.Exchange) <-chan domain.Exchange {
	var out = make(chan domain.Exchange)
	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, in := range ins {
		go func(in <-chan domain.Exchange) {
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

func CreateHashTable(exchanges, pairNames []string) map[string]map[string]map[string]float64 {
	var m = make(map[string]map[string]map[string]float64, 3)

	for i := range exchanges {
		m[exchanges[i]] = make(map[string]map[string]float64, 5)
		for j := range pairNames {
			m[exchanges[i]][pairNames[j]] = make(map[string]float64, 4)
			m[exchanges[i]][pairNames[j]]["min"] = 0
		}
	}

	return m
}

func (exchangeRepo *exchangeRepository) Aggregate(exchanges, pairNames []string, m map[string]map[string]map[string]float64) []domain.Exchange {
	exchangesData, err := exchangeRepo.redisRepository.ReadAll(exchanges, pairNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	for i := range exchangesData {
		var data = exchangesData[i]

		m[data.Exchange][data.Symbol]["avg"] += data.Price

		if m[data.Exchange][data.Symbol]["min"] > data.Price || m[data.Exchange][data.Symbol]["min"] == 0 {
			m[data.Exchange][data.Symbol]["min"] = data.Price
		}
		if m[data.Exchange][data.Symbol]["max"] < data.Price {
			m[data.Exchange][data.Symbol]["max"] = data.Price
		}

		m[data.Exchange][data.Symbol]["count"]++
	}

	return exchangesData
}

func (exchangeRepo *exchangeRepository) GetAggregatedData(exchanges, pairNames []string, m map[string]map[string]map[string]float64) []domain.Exchanges {
	var aggregatedData []domain.Exchanges

	for i := range exchanges {
		for j := range pairNames {
			var exchange = domain.Exchanges{
				PairName:  pairNames[j],
				Exchange:  exchanges[i],
				Timestamp: time.Now(),
				AvgPrice:  m[exchanges[i]][pairNames[j]]["avg"] / m[exchanges[i]][pairNames[j]]["count"],
				MinPrice:  m[exchanges[i]][pairNames[j]]["min"],
				MaxPrice:  m[exchanges[i]][pairNames[j]]["max"],
			}

			aggregatedData = append(aggregatedData, exchange)
		}
	}

	return aggregatedData
}

func (exchangeRepo *exchangeRepository) WriteToStorage(exchanges []string, ticker *time.Ticker, done chan bool) {
	var pairNames = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var m = CreateHashTable(exchanges, pairNames)
				var exchangesData = exchangeRepo.Aggregate(exchanges, pairNames, m)
				var aggregatedData = exchangeRepo.GetAggregatedData(exchanges, pairNames, m)

				var err = exchangeRepo.postgresRepository.Write(aggregatedData)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				err = exchangeRepo.redisRepository.DeleteAll(exchangesData)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			}
		}
	}()
}

func (exchangeRepo *exchangeRepository) LiveMode() {
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}
	var exchangesAddress = []string{"172.22.0.5:40101", "172.22.0.6:40102", "172.22.0.7:40103"}
	var outSlice = make([]<-chan domain.Exchange, 15)
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
		var err = exchangeRepo.redisRepository.Write(i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	done <- true
	ticker.Stop()
}

func Generator() <-chan string {
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

func (exchangeRepo *exchangeRepository) TestMode() {
	var exchanges = []string{"exchange_test"}
	var outSlice = make([]<-chan domain.Exchange, 10)
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
