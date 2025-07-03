package exchangerepository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"marketflow/internal/core/ports"
	"net"
	"os"
	"strings"
	"sync"
)

type exchangeRepository struct {
	redisRepository ports.RedisRepository
}

func NewExchangeRepository(redisRepo ports.RedisRepository) *exchangeRepository {
	return &exchangeRepository{
		redisRepository: redisRepo,
	}
}

type Exchange struct {
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func Distributor(conn net.Conn) <-chan string {
	var scanner = bufio.NewScanner(conn)
	var out = make(chan string)

	go func() {
		defer close(out)

		for scanner.Scan() {
			var result = scanner.Text()

			out <- result
		}
	}()

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

func (exchangeRepo *exchangeRepository) GetData() {
	conn1, err := net.Dial("tcp", ":40101")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer conn1.Close()

	conn2, err := net.Dial("tcp", ":40102")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer conn2.Close()

	conn3, err := net.Dial("tcp", ":40103")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer conn3.Close()

	var out1 = Distributor(conn1)
	var out2 = Distributor(conn2)
	var out3 = Distributor(conn3)

	var outSlice = make([]<-chan string, 15)

	for i := 0; i < 5; i++ {
		outSlice[i] = Worker(out1, "exchange1")
	}

	for i := 5; i < 10; i++ {
		outSlice[i] = Worker(out2, "exchange2")
	}

	for i := 10; i < 15; i++ {
		outSlice[i] = Worker(out3, "exchange3")
	}

	var merged = Merger(outSlice...)

	exchangeRepo.redisRepository.Write(merged)
}
