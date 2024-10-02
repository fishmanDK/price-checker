// package service

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/fatih/color"
// 	"github.com/fishmanDK/price_checker/internal/storage"
// 	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
// )

// type Service struct {
// 	httpClient *http.Client
// 	storage    *storage.Storage
// }

// func NewService(_storage *storage.Storage) *Service {
// 	client := &http.Client{}
// 	return &Service{
// 		storage:    _storage,
// 		httpClient: client,
// 	}
// }

// func (s *Service) Start() {
// 	fileData, err := os.Create("data/data.txt")
// 	if err != nil {
// 		log.Fatalf("Не удалось создать файл данных: %v", err)
// 	}
// 	defer fileData.Close()

// 	for {
// 		startTime := time.Now()
// 		nextRun := startTime.Truncate(5 * time.Minute).Add(5 * time.Minute)
// 		time.Sleep(time.Until(nextRun))

// 		var wg sync.WaitGroup
// 		for symb := range tokensTrade {
// 			wg.Add(1)
// 			go func(symb string) {
// 				defer wg.Done()
// 				s.CollectInfo(symb, fileData)
// 			}(symb)
// 		}
// 		wg.Wait()
// 	}
// }

// func (s *Service) getPriceToken(symb string) (float64, error) {
// 	query := "https://api.mexc.com/api/v3/avgPrice?symbol=" + symb + "USDT"
// 	resp, err := s.httpClient.Get(query)
// 	if err != nil {
// 		return 0.0, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.Status[:3] != "200"{
// 		if resp.Status[:3] == "429"{
// 			red.Println("429")
// 			log.Fatal()
// 		}
// 		red.Println(fmt.Sprintf("%s - %s", symb, resp.Status))
// 	}
// 	var res SymbolPriceToken
// 	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
// 		return 0.0, err
// 	}

// 	return strconv.ParseFloat(res.Price, 64)
// }

// func (s *Service) CollectInfo(symb string, fileData *os.File) {
// 	priceToken, err := s.getPriceToken(symb)
// 	if err != nil {
// 		log.Println(symb, err)
// 		return
// 	}

// 	org := "price_checker"
// 	bucket := "tokens"
// 	writeAPI := s.storage.Client.WriteAPIBlocking(org, bucket)

// 	now := time.Now().In(time.FixedZone("MSK", 3*3600)).Truncate(5 * time.Minute)
// 	tags := map[string]string{"symbol": symb}
// 	fields := map[string]interface{}{"Price": priceToken}

// 	point := influxdb2.NewPoint(symb+"/USDT", tags, fields, now)
// 	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
// 		log.Fatal(err)
// 	}

// 	prices := s.fetchHistoricalPrices(symb)
// 	if prices == nil {
// 		return
// 	}

// 	priceRatios := calculatePriceRatios(prices, priceToken)

// 	s.printResults(symb, prices, priceRatios)
// }

// func (s *Service) fetchHistoricalPrices(symb string) []float64 {
// 	queryAPI := s.storage.Client.QueryAPI("price_checker")
// 	var wg sync.WaitGroup
// 	times := []time.Time{
// 		time.Now().Add(-5 * time.Minute).Truncate(time.Minute),
// 		time.Now().Add(-15 * time.Minute).Truncate(time.Minute),
// 		time.Now().Add(-30 * time.Minute).Truncate(time.Minute),
// 		time.Now().Add(-1 * time.Hour).Truncate(time.Minute),
// 		time.Now().Add(-4 * time.Hour).Truncate(time.Minute),
// 		time.Now().Add(-24 * time.Hour).Truncate(time.Minute),
// 	}

// 	prices := make([]float64, len(times))
// 	for idx, date := range times {
// 		wg.Add(1)
// 		go func(idx int, date time.Time) {
// 			defer wg.Done()
// 			// query := fmt.Sprintf(`
// 			// 	from(bucket: "tokens")
// 			// 	|> range(start: %s, stop: %s)
// 			// 	|> filter(fn: (r) => r._measurement == "%s/USDT" and r._value > 0)`,
// 			// 	date.Format(time.RFC3339), date.Add(1*time.Minute).Format(time.RFC3339), symb)

// 			query := fmt.Sprintf(`from(bucket: "%s")
// 			|> range(start: %s, stop: %s)
// 			|> filter(fn: (r) => r._measurement == "%s/USDT" and r._value > 0)`,
// 				bucket, date.Format(time.RFC3339), date.Add(1*time.Minute).Format(time.RFC3339), symb)

// 			results, err := queryAPI.Query(context.Background(), query)
// 			if err != nil {
// 				log.Printf("Ошибка запроса для %s: %v", symb, err)
// 				return
// 			}

// 			if results.Next() {
// 				prices[idx] = results.Record().Value().(float64)
// 			} else {
// 				prices[idx] = 0.0
// 			}
// 			if err := results.Err(); err != nil {
// 				log.Printf("Ошибка результата для %s: %v", symb, err)
// 			}
// 		}(idx, date)
// 	}
// 	wg.Wait()
// 	return prices
// }

// func calculatePriceRatios(prices []float64, currentPrice float64) []float64 {
// 	priceRatios := make([]float64, len(prices))
// 	for i, price := range prices {
// 		if price > 0 {
// 			priceRatios[i] = ((price/currentPrice) - 1) * 100
// 		} else {
// 			priceRatios[i] = 0.0
// 		}
// 	}
// 	return priceRatios
// }

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	kafkaproducer "github.com/fishmanDK/price_checker/internal/modules/kafkaProducer"
	"github.com/fishmanDK/price_checker/internal/storage"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Service struct {
	httpClient *http.Client
	storage    *storage.Storage
	kafkaProducer *kafkaproducer.KafkaProducer
}

func NewService(_storage *storage.Storage) *Service {
	client := &http.Client{}
	_kafkaProducer, err := kafkaproducer.NewKafkaProducer("localhost:29092", "ratios-info")
	if err != nil{
		log.Fatal(err.Error())
	}

	return &Service{
		storage:    _storage,
		httpClient: client,
		kafkaProducer: _kafkaProducer,
	}
}

func (s *Service) Start() {
	fileData, err := os.Create("data/data.txt")
	if err != nil {
		log.Fatalf("Не удалось создать файл данных: %v", err)
	}
	defer fileData.Close()

	for {
		startTime := time.Now()
		nextRun := startTime.Truncate(5 * time.Minute).Add(5 * time.Minute)
		time.Sleep(time.Until(nextRun))

		var wg sync.WaitGroup
		for symb := range tokensTrade {
			wg.Add(1)
			go func(symb string) {
				defer wg.Done()
				s.CollectInfo(symb, fileData)
			}(symb)
		}
		wg.Wait()
	}
}

func (s *Service) getPriceToken(symb string) (float64, error) {
	query := "https://api.mexc.com/api/v3/avgPrice?symbol=" + symb + "USDT"
	resp, err := s.httpClient.Get(query)
	if err != nil {
		log.Printf("Ошибка при выполнении HTTP-запроса к API МЭКС: %v", err)
		return 0.0, err
	}
	defer resp.Body.Close()

	if resp.Status[:3] != "200" {
		if resp.Status[:3] == "429" {
			red.Println("429 - Rate limit exceeded")
			log.Fatal("Rate limit exceeded")
		}
		red.Println(fmt.Sprintf("%s - %s", symb, resp.Status))
		log.Printf("Неверный статус ответа для символа %s: %s", symb, resp.Status)
	}

	var res SymbolPriceToken
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Printf("Ошибка при декодировании JSON для символа %s: %v", symb, err)
		return 0.0, err
	}

	price, err := strconv.ParseFloat(res.Price, 64)
	if err != nil {
		log.Printf("Ошибка при парсинге цены для символа %s: %v", symb, err)
		return 0.0, err
	}

	return price, nil
}

func (s *Service) CollectInfo(symb string, fileData *os.File) {
	priceToken, err := s.getPriceToken(symb)
	if err != nil {
		log.Printf("Ошибка при получении цены токена для символа %s: %v", symb, err)
		return
	}

	org := "price_checker"
	bucket := "tokens"
	writeAPI := s.storage.Client.WriteAPIBlocking(org, bucket)

	now := time.Now().In(time.FixedZone("MSK", 3*3600)).Truncate(5 * time.Minute)
	tags := map[string]string{"symbol": symb}
	fields := map[string]interface{}{"Price": priceToken}

	point := influxdb2.NewPoint(symb+"/USDT", tags, fields, now)
	err = writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		log.Printf("Ошибка при записи данных в InfluxDB для символа %s: %v", symb, err)
		return
	}

	prices := s.fetchHistoricalPrices(symb)
	if prices == nil {
		log.Printf("Не удалось получить исторические цены для символа %s", symb)
		return
	}

	priceRatios := calculatePriceRatios(prices, priceToken)

	s.printResults(symb, prices, priceRatios)
}

func (s *Service) fetchHistoricalPrices(symb string) []float64 {
	queryAPI := s.storage.Client.QueryAPI("price_checker")
	var wg sync.WaitGroup
	times := []time.Time{
		time.Now().Add(-5 * time.Minute).Truncate(time.Minute),
		time.Now().Add(-15 * time.Minute).Truncate(time.Minute),
		time.Now().Add(-30 * time.Minute).Truncate(time.Minute),
		time.Now().Add(-1 * time.Hour).Truncate(time.Minute),
		time.Now().Add(-4 * time.Hour).Truncate(time.Minute),
		time.Now().Add(-24 * time.Hour).Truncate(time.Minute),
	}

	prices := make([]float64, len(times))
	for idx, date := range times {
		wg.Add(1)
		go func(idx int, date time.Time) {
			defer wg.Done()
			// query := fmt.Sprintf(`
			// 	from(bucket: "tokens")
			// 	|> range(start: %s, stop: %s)
			// 	|> filter(fn: (r) => r._measurement == "%s/USDT" and r._value > 0)`,
			// 	date.Format(time.RFC3339), date.Add(1*time.Minute).Format(time.RFC3339), symb)
			query := fmt.Sprintf(`from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r._measurement == "%s/USDT" and r._value > 0)`,
				bucket, date.Format(time.RFC3339), date.Add(1*time.Minute).Format(time.RFC3339), symb)
			results, err := queryAPI.Query(context.Background(), query)
			if err != nil {
				log.Printf("Ошибка запроса для %s: %v", symb, err)
				return
			}
			if results.Next() {
				prices[idx] = results.Record().Value().(float64)
			} else {
				prices[idx] = 0.0
			}
			if err := results.Err(); err != nil {
				log.Printf("Ошибка результата для %s: %v", symb, err)
			}
		}(idx, date)
	}
	wg.Wait()
	return prices
}

func calculatePriceRatios(prices []float64, currentPrice float64) []float64 {
	priceRatios := make([]float64, len(prices))
	for i, price := range prices {
		if price > 0 {
			priceRatios[i] = ((price / currentPrice) - 1) * 100
		} else {
			priceRatios[i] = 0.0
		}
	}
	return priceRatios
}

var durations = []string{"5 мин", "15 мин", "30 мин", "1 ч", "4 ч", "1 д"}
var settings = []struct {
	duration   string
	diffGrowth float64
	diffDrop   float64
}{
	{"5 мин", 0.8, -0.8},
	{"15 мин", 1.5, -1.5},
	{"30 мин", 2.5, -2.5},
	{"1 ч", 6.0, -6.0},
	{"4 ч", 8.8, -8.8},
	{"1 д", 15.0, -15.0},
}

var red = color.New(color.FgRed, color.Bold)

func (s *Service) printResults(symb string, prices, priceRatios []float64) {
	for i, set := range settings {
		if priceRatios[i] >= set.diffGrowth || priceRatios[i] <= set.diffDrop {
			red.Printf("%s: %s - %f\n", set.duration, symb, priceRatios[i])

			err := s.kafkaProducer.PublishMessage(symb, priceRatios[i], set.duration)
			if err != nil {
				log.Printf("Error publishing message to Kafka: %v", err)
			}
			break
		}
	}
}
