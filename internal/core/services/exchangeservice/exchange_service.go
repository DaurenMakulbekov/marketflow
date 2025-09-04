package exchangeservice

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
	storage            ports.Storage
	live               bool
	liveStarted        bool
	test               bool
}

func NewExchangeService(exchangeRepo ports.ExchangeRepository, redisRepo ports.RedisRepository, postgresRepo ports.PostgresRepository, storageRepo ports.Storage) *exchangeService {
	return &exchangeService{
		exchangeRepository: exchangeRepo,
		redisRepository:    redisRepo,
		postgresRepository: postgresRepo,
		storage:            storageRepo,
	}
}

func (exchangeServ *exchangeService) Distributor(exchanges []string) []<-chan domain.Exchange {
	var outSlice = make([]<-chan domain.Exchange, 15)
	var index int = 0

	for i := range exchanges {
		var out <-chan string

		if exchangeServ.liveStarted == false && exchangeServ.live == true {
			out = exchangeServ.exchangeRepository.GetFromExchange(exchanges[i])
		} else {
			out = exchangeServ.exchangeRepository.Generator()
		}

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

			var id = strconv.FormatInt(result.Timestamp, 10)
			var result1 = domain.Exchange{
				ID:        id,
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

func (exchangeServ *exchangeService) Aggregate(exchanges, pairNames []string, m map[string]map[string]map[string]float64) ([]domain.Exchange, []domain.Exchange) {
	var result []domain.Exchange

	exchangesData, err := exchangeServ.redisRepository.ReadAll(exchanges, pairNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if len(exchangesData) > 0 {
		result = append(result, exchangesData...)
	}

	var storageData = exchangeServ.storage.GetAll()
	if len(storageData) > 0 {
		result = append(result, storageData...)
	}

	for i := range result {
		var data = result[i]

		m[data.Exchange][data.Symbol]["avg"] += data.Price

		if m[data.Exchange][data.Symbol]["min"] > data.Price || m[data.Exchange][data.Symbol]["min"] == 0 {
			m[data.Exchange][data.Symbol]["min"] = data.Price
		}
		if m[data.Exchange][data.Symbol]["max"] < data.Price {
			m[data.Exchange][data.Symbol]["max"] = data.Price
		}

		m[data.Exchange][data.Symbol]["count"]++
	}

	return exchangesData, storageData
}

func GetAggregatedData(exchanges, pairNames []string, m map[string]map[string]map[string]float64) []domain.Exchanges {
	var aggregatedData []domain.Exchanges

	for i := range exchanges {
		for j := range pairNames {
			if m[exchanges[i]][pairNames[j]]["count"] > 0 {
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
	}

	return aggregatedData
}

func (exchangeServ *exchangeService) WriteToStorage(exchanges, pairNames []string, ticker *time.Ticker, done chan bool) {
	var storage []domain.Exchanges
	var cache []domain.Exchange
	var localstorage []domain.Exchange

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if len(localstorage) > 0 {
					exchangeServ.storage.DeleteAll(localstorage)
					localstorage = []domain.Exchange{}
				}
				if len(cache) > 0 {
					var err = exchangeServ.redisRepository.DeleteAll(cache)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
					} else {
						cache = []domain.Exchange{}
					}
				}

				var m = CreateHashTable(exchanges, pairNames)
				cacheData, localstorageData := exchangeServ.Aggregate(exchanges, pairNames, m)
				var aggregatedData = GetAggregatedData(exchanges, pairNames, m)
				storage = append(storage, aggregatedData...)

				var err = exchangeServ.postgresRepository.Write(storage)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					storage = []domain.Exchanges{}
				}

				cache = append(cache, cacheData...)
				localstorage = append(localstorage, localstorageData...)
			}
		}
	}()
}

func (exchangeServ *exchangeService) RedisConnect(doneRedisConn chan bool, healthy *bool) {
	var err = exchangeServ.redisRepository.CheckConnection()
	if err != nil {
		log.Printf("Error: %v", err)

		var ticker = time.NewTicker(time.Second)
		var done = make(chan bool)

		go func() {
			defer close(done)

			for {
				select {
				case <-doneRedisConn:
					return
				case <-ticker.C:
					exchangeServ.redisRepository.Reconnect()
					var err = exchangeServ.redisRepository.CheckConnection()
					if err == nil {
						log.Println("Connected to Redis")
						return
					}
				}
			}
		}()

		<-done
		ticker.Stop()
	}
	*healthy = true
}

func (exchangeServ *exchangeService) RunLive() {
	exchanges, pairNames := exchangeServ.exchangeRepository.GetExchanges()
	var out = exchangeServ.Distributor(exchanges)
	exchangeServ.liveStarted = true
	var merged = Merger(out...)

	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)
	defer close(done)

	exchangeServ.WriteToStorage(exchanges, pairNames, ticker, done)

	var healthy bool = true
	var doneRedisConn = make(chan bool)
	defer close(doneRedisConn)

	for i := range merged {
		if healthy {
			var err = exchangeServ.redisRepository.Write(i)
			if err != nil {
				healthy = false
				exchangeServ.storage.Write(i)
				go exchangeServ.RedisConnect(doneRedisConn, &healthy)

				fmt.Fprintln(os.Stderr, err.Error())
			}
		} else {
			exchangeServ.storage.Write(i)
		}
	}

	doneRedisConn <- true
	done <- true
	ticker.Stop()
}

func (exchangeServ *exchangeService) LiveMode() {
	if exchangeServ.live == true {
		return
	}
	if exchangeServ.test == true {
		exchangeServ.exchangeRepository.CloseTest()
	}
	exchangeServ.live = true
	exchangeServ.test = false

	if exchangeServ.liveStarted == false {
		go exchangeServ.RunLive()
	}
}

func (exchangeServ *exchangeService) RunTest() {
	exchanges, pairNames := exchangeServ.exchangeRepository.GetExchangesTest()
	var out = exchangeServ.Distributor(exchanges)
	var merged = Merger(out...)

	var ticker = time.NewTicker(60 * time.Second)
	var done = make(chan bool)
	defer close(done)

	exchangeServ.WriteToStorage(exchanges, pairNames, ticker, done)

	for i := range merged {
		var err = exchangeServ.redisRepository.Write(i)
		if err != nil {
			exchangeServ.storage.Write(i)

			//fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	done <- true
	ticker.Stop()
}

func (exchangeServ *exchangeService) TestMode() {
	if exchangeServ.test == true {
		return
	}
	exchangeServ.test = true
	exchangeServ.live = false

	go exchangeServ.RunTest()
}
