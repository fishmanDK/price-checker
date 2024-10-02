// package modules

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// type Result struct {
// 	RetCode     int      `json:"retCode"`
// 	RetMsg      string   `json:"retMsg"`
// 	Result       ResultItem `json:"result"`
// 	RetExtInfo  map[string]interface{} `json:"retExtInfo"`
// 	Time         int64 `json:"time"`
// }

// // ResultItem структура, представляющая основной результат запроса
// type ResultItem struct {
// 	Symbol       string `json:"symbol"`
// 	Category     string `json:"category"`
// 	List         []OpenInterestData `json:"list"`
// 	NextPageCursor string `json:"nextPageCursor"`
// }

// // OpenInterestData структура, представляющая данные о открытом интересе
// type OpenInterestData struct {
// 	OpenInterest string `json:"openInterest"`
// 	Timestamp   string `json:"timestamp"` // Изменили тип на string
// }

// func GetOpenInterest(){
// 	url := "https://api-testnet.bybit.com/v5/market/open-interest?category=inverse&symbol=BTCUSDT&intervalTime=5min"
// 	method := "GET"

// 	client := &http.Client {}
// 	req, err := http.NewRequest(method, url, nil)

// 	if err != nil {
// 	  fmt.Println(err)
// 	  return
// 	}
// 	res, err := client.Do(req)
// 	if err != nil {
// 	  fmt.Println(err)
// 	  return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 	  fmt.Println(err)
// 	  return
// 	}
// 	var result Result
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		fmt.Println("Ошибка при разборе JSON:", err)
// 		return
// 	}

//		fmt.Println(result)
//	}

package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Result struct {
	RetCode     int          `json:"retCode"`
	RetMsg      string       `json:"retMsg"`
	Result      ResultItem   `json:"result"`
	RetExtInfo  interface{}  `json:"retExtInfo"` // Changed to interface{} for flexibility
	Time        int64        `json:"time"`
}

type ResultItem struct {
	Symbol        string              `json:"symbol"`
	Category      string              `json:"category"`
	List          []OpenInterestData   `json:"list"`
	NextPageCursor string             `json:"nextPageCursor"`
}

type OpenInterestData struct {
	OpenInterest string `json:"openInterest"`
	Timestamp    string `json:"timestamp"`
}

func GetOpenInterest() {
	url := "https://api-testnet.bybit.com/v5/market/open-interest?category=inverse&symbol=BTCUSDT&intervalTime=1h"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Error: Received non-200 response status:", res.Status)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if result.RetCode != 0 {
		fmt.Printf("API Error: %s (Code: %d)\n", result.RetMsg, result.RetCode)
		return
	}

	// Create graph
	createGraph(result.Result.List)
	createGraphVisual(result.Result.List)
}

func createGraph(data []OpenInterestData) {
	if len(data) == 0 {
		fmt.Println("No data available to create graph.")
		return
	}

	// Clear screen (Note: This may not work in all environments)
	fmt.Print("\033[H\033[2J")

	// Header of the graph
	fmt.Println("График открытого интереса BTCUSDT")
	fmt.Println("Время в MSK")

	// Create slices to hold parsed open interest values and formatted timestamps
	values := make([]float64, len(data))
	timestamps := make([]string, len(data))

	for i, item := range data {
		value, err := strconv.ParseFloat(item.OpenInterest, 64)
		if err != nil {
			fmt.Printf("Error parsing open interest for timestamp %s: %v\n", item.Timestamp, err)
			continue // Skip this entry if there's an error
		}
		values[i] = value

		// Convert timestamp from milliseconds to seconds
		msTimestamp, err := strconv.ParseInt(item.Timestamp, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing timestamp: %v\n", err)
			continue // Skip this entry if there's an error
		}

		// Create a time.Time object and convert it to Moscow time
		t := time.Unix(msTimestamp/1000, 0).In(time.FixedZone("MSK", 3*3600)) // MSK is UTC+3
		timestamps[i] = t.Format("2006-01-02 15:04:05") // Format the time as needed
	}

	minValue := values[0]
	maxValue := values[0]

	for _, v := range values {
		if v < minValue {
			minValue = v
		}
		if v > maxValue {
			maxValue = v
		}
	}

	fmt.Printf("\n%-20s | %-15s\n", "Timestamp", "Open Interest")
	fmt.Println(strings.Repeat("-", 40))

	for i := range values {
		if i < len(timestamps) {
			fmt.Printf("%-20s | %-15.2f\n", timestamps[i], values[i])
		}
	}

	fmt.Printf("\nМинимальное значение: %.2f\n", minValue)
	fmt.Printf("Максимальное значение: %.2f\n", maxValue)
}

func createGraphVisual(data []OpenInterestData) {
	if len(data) == 0 {
		fmt.Println("Нет данных для создания графика.")
		return
	}

	values := make([]opts.LineData, len(data))
	timestamps := make([]string, len(data))

	for i, item := range data {
		value, err := strconv.ParseFloat(item.OpenInterest, 64)
		if err != nil {
			fmt.Printf("Ошибка при парсинге открытого интереса для временной метки %s: %v\n", item.Timestamp, err)
			continue
		}

		msTimestamp, err := strconv.ParseInt(item.Timestamp, 10, 64)
		if err != nil {
			fmt.Printf("Ошибка при парсинге временной метки: %v\n", err)
			continue
		}

		t := time.Unix(msTimestamp/1000, 0).In(time.FixedZone("MSK", 3*3600))
		timestamps[i] = t.Format("2006-01-02 15:04:05")
		values[i] = opts.LineData{Value: value}
	}
	
	reversedTimestamps := make([]string, len(data))
	reversedValues := make([]opts.LineData, len(data))

	maxV, minV := 0.0, math.MaxFloat64
	for i := range data {
		f, _ := strconv.ParseFloat(data[i].OpenInterest, 64)
		if maxV < f{maxV = f}
		if minV > f{minV = f}

		reversedTimestamps[i] = timestamps[len(data)-1-i]
		reversedValues[i] = values[len(data)-1-i]
	}

	hWindow := 0.05
	window := (((maxV + minV) / 2) / 100) * hWindow

	const actionWithEchartsInstance = `
	var myChart = echarts.init(document.getElementById('container'), null, {
    width: 6000,
    height: 4000
  });
    `

	line := charts.NewLine()

	line.AddJSFuncs(actionWithEchartsInstance)

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "График открытого интереса BTCUSDT",
			Subtitle: "Данные по открытому интересу за период времени",
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithYAxisOpts(opts.YAxis{
            Name: "Открытый интерес",
            Max: int(((maxV + minV) / 2) + window) ,
			Min: int(((maxV + minV) / 2) - window) ,
        }),
	)

	line.SetXAxis(reversedTimestamps).
		AddSeries("Open Interest", reversedValues,
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
			
		)

	f, err := os.Create("open_interest_chart.html")
	if err != nil {
	    fmt.Println("Ошибка при создании HTML файла:", err)
	    return
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
	    fmt.Println("Ошибка при рендеринге графика:", err)
	    return
	}

	fmt.Println("График успешно создан: open_interest_chart.html")
}