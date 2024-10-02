package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func getHistoricalData() {
	url := "https://www.coingecko.com/price_charts/export/1/usd.xls"
	fileName := "btc-usd-max.xls"

	// Создаем новый HTTP-запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Устанавливаем необходимые заголовки
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "ru,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", "https://www.coingecko.com/en/coins/bitcoin/historical_data?start=2020-09-22&end=2024-09-21")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 YaBrowser/24.7.0.0 Safari/537.36")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: статус ответа", resp.Status)
		return
	}

	// Создаем файл для сохранения
	out, err := os.Create("./data/"+fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer out.Close()

	// Копируем данные из ответа в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	fmt.Println("Файл успешно загружен и сохранен как", fileName)
}