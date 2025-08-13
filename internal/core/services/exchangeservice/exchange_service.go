package exchangeservice

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"marketflow/internal/core/domain"
	"marketflow/internal/core/ports"
)

type exchangeService struct {
	exchangeRepository ports.ExchangeRepository
	redisRepository    ports.RedisRepository
	postgresRepository ports.PostgresRepository
}

func NewExchangeService(exchangeRepo ports.ExchangeRepository, redisRepo ports.RedisRepository, postgresRepo ports.PostgresRepository) *exchangeService {
	return &exchangeService{
		exchangeRepository: exchangeRepo,
		redisRepository:    redisRepo,
		postgresRepository: postgresRepo,
	}
}

func (exchangeServ *exchangeService) Distributor(exchangesAddress []string, exchanges []string) []<-chan domain.Exchange {
	var outSlice = make([]<-chan domain.Exchange, 15)
	var index int = 0

	for i := range exchangesAddress {
		var out = exchangeServ.exchangeRepository.GetFromExchange(exchangesAddress[i])

		for j := 0; j < 5; j++ {
			outSlice[index] = Worker(out, exchanges[i])
			index++
		}
	}

	return outSlice
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

func (exchangeServ *exchangeService) Aggregate(exchanges, pairNames []string, m map[string]map[string]map[string]float64) []domain.Exchange {
	exchangesData, err := exchangeServ.redisRepository.ReadAll(exchanges, pairNames)
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

func GetAggregatedData(exchanges, pairNames []string, m map[string]map[string]map[string]float64) []domain.Exchanges {
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

func (exchangeServ *exchangeService) WriteToStorage(exchanges []string, ticker *time.Ticker, done chan bool) {
	var pairNames = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var m = CreateHashTable(exchanges, pairNames)
				var exchangesData = exchangeServ.Aggregate(exchanges, pairNames, m)
				var aggregatedData = GetAggregatedData(exchanges, pairNames, m)

				var err = exchangeServ.postgresRepository.Write(aggregatedData)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				err = exchangeServ.redisRepository.DeleteAll(exchangesData)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			}
		}
	}()
}

func (exchangeServ *exchangeService) LiveMode() {
	var exchanges = []string{"exchange1", "exchange2", "exchange3"}
	var exchangesAddress = []string{"172.22.0.5:40101", "172.22.0.6:40102", "172.22.0.7:40103"}

	var out = exchangeServ.Distributor(exchangesAddress, exchanges)

	var merged = Merger(out...)
	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)

	exchangeServ.WriteToStorage(exchanges, ticker, done)

	for i := range merged {
		var err = exchangeServ.redisRepository.Write(i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	done <- true
	ticker.Stop()
}

func (exchangeServ *exchangeService) TestMode() {
	var exchanges = []string{"exchange_test"}
	var outSlice = make([]<-chan domain.Exchange, 10)
	var out = exchangeServ.exchangeRepository.Generator()

	for i := 0; i < 5; i++ {
		outSlice[i] = Worker(out, "exchange_test")
	}

	var merged = Merger(outSlice...)
	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)

	exchangeServ.WriteToStorage(exchanges, ticker, done)

	for i := range merged {
		exchangeServ.redisRepository.Write(i)
	}

	done <- true
	ticker.Stop()
}
