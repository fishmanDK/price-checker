package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fishmanDK/price_checker/internal/service"
	"github.com/fishmanDK/price_checker/internal/storage"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type TickerPrice []SymbolPriceToken

type SymbolPriceToken struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

var dct = map[string]string{}

const (
	X    = 3
	stop = 0.05
)

type TokenPositionInfo struct {
	ProcurementPrice float64
	PlusExitPrice    float64
	MinusExitPrice   float64
}

var wallet = 100.0
var marja = 0.0

var importantTokens = map[string]bool{
	"BTC":  true,
	"ETH":  true,
	"ADA":  true,
	"BNB":  true,
	"SOL":  true,
	"USDC": true,
	"XRP":  true,
	"DOGE": true,
	"TON":  true,
	"TRX":  true,
}

const urlInflux = "http://localhost:8086"

// func main(){
// 	modules.GetOpenInterest()
// 	keyboard.Open()
// }

func main() {
	// deleteData() 
	_storage := storage.NewStorage(urlInflux)
	_service := service.NewService(_storage)
	_service.Start()
	// start := time.Now()
	// fileData, _ := os.Create("data.txt")
	// errorsFile, _ := os.Open("errors.txt")
	// defer fileData.Close()
	// red := color.New(color.FgRed, color.Bold)
	// green := color.New(color.FgGreen, color.Bold)

	//
	// _storage := storage.NewStorage(urlInflux)

	// org := "price_checker"
	// bucket := "tokens"
	// writeAPI := _storage.WriteAPIBlocking(org, bucket)

	// client := http.Client{}

	// wg := &sync.WaitGroup{}

	// workerStart := func(symb string, wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	query := "https://api.mexc.com/api/v3/avgPrice?symbol=" + symb + "USDT"
	// 	resp, err := client.Get(query)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		return
	// 	}

	// 	if resp.Status == "403" {
	// 		time.Sleep(time.Second)
	// 		resp, err = client.Get(query)
	// 		if err != nil {
	// 			log.Println(err.Error())
	// 			return
	// 		}

	// 		if resp.Status != "200" {
	// 			errorsFile.WriteString(fmt.Sprintf("%s\n", symb))
	// 		}
	// 	}

	// 	if importantTokens[symb] {
	// 		red.Println(symb, "-", resp.Status)
	// 	} else {
	// 		fmt.Println(symb, "-", resp.Status)
	// 	}

	// 	var res SymbolPriceToken
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Println("Ошибка чтения тела запроса")
	// 		return
	// 	}
	// 	defer resp.Body.Close()
	// 	if err := json.Unmarshal(body, &res); err != nil {
	// 		log.Println("Ошибка обработки JSON")
	// 		return
	// 	}

	// 	f, _ := strconv.ParseFloat(res.Price, 64)

	// 	tokensTrade[symb].Prices1H[0] = f
	// 	tokensTrade[symb].PricesAlong1H.Current += 1
	// 	tokensTrade[symb].PricesAlong1H.TotalSum += f

	// 	tokensTrade[symb].Prices30Min[0] = f
	// 	tokensTrade[symb].PricesAlont30Min.Current += 1
	// 	tokensTrade[symb].PricesAlont30Min.TotalSum += f

	// 	tokensTrade[symb].PriceStart = f
	// 	tokensTrade[symb].LastPrice = f
	// 	tokensTrade[symb].Start = time.Now()

	// 	now := time.Now()

	// 	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	// 	if err != nil {
	// 		fmt.Println("Ошибка при загрузке часового пояса:", err)
	// 		return
	// 	}

	// 	moscowTime := now.In(moscowLocation)

	// 	tags := map[string]string{
	// 		"symbol": symb,
	// 	}
	// 	fields := map[string]interface{}{
	// 		"Price": f,
	// 	}

	// 	point := influxdb2.NewPoint(symb+"/USDT", tags, fields, moscowTime)

	// 	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// for symb := range tokensTrade {
	// 	wg.Add(1)
	// 	go workerStart(symb, wg)
	// }
	// wg.Wait()

	// fmt.Println(time.Since(start))

	// for {
	// 	time.Sleep(time.Minute * 5)
	// 	fmt.Println()
	// 	for symb := range tokensTrade {
	// 		go worker(&client, symb, red, green, fileData)
	// 	}
	// }
}


func deleteData() {
	url := "http://localhost:8086" // URL вашего InfluxDB
	token := "VRsOE9gM1OcQS5VFEwiR3pZERheEi4fXm909k_C7aaALGql509v3nNH4oWe5ziW70K0bXDG-P3Js-YxJfndHFA=="           // Ваш токен доступа
	org := "price_checker"               // Ваша организация
	bucket := "tokens"
         // Ваш бакет

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	// Удаляем все данные из бакета
	start := "1970-01-01T00:00:00Z" // Начало временного диапазона
	stop := "2025-01-01T00:00:00Z" 
	// Создаем тело запроса
	body := map[string]interface{}{
		"start": start,
		"stop":  stop,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Ошибка при маршализации JSON: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/delete?org=%s&bucket=%s", url, org, bucket), bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")

	clientH := &http.Client{}
	resp, err := clientH.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Fatalf("Ошибка при удалении данных. Статус: %s", resp.Status)
	}

	log.Println("Все данные успешно удалены из бакета:", bucket)
}