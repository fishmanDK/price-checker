package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fishmanDK/price_checker/internal/core"
	kafkaClient "github.com/fishmanDK/price_checker/internal/kafka"
	"github.com/fishmanDK/price_checker/internal/logger"
	"github.com/fishmanDK/price_checker/internal/marketdata/storage"
	kafkaMessages "github.com/fishmanDK/price_checker/proto/kafka"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Checker struct {
		httpClient    *http.Client
		storage       storage.Storage
		kafkaClientProducer kafkaClient.Producer
		log           logger.Logger
	}

	Config struct {
	}
)

func NewConfig() *Config { return &Config{} } //TODO

func NewChecker(storage storage.Storage, kafkaProducer kafkaClient.Producer, log logger.Logger, cfg *Config) *Checker {
	client := &http.Client{}

	return &Checker{
		httpClient:    client,
		storage:       storage,
		kafkaClientProducer: kafkaProducer,
		log:           log,
	}
}

func (c *Checker) Start() {
	for {
		tokensTrade, err := c.storage.GetTokensSymbols()
		if err != nil {
			c.log.Error(errors.Wrap(err, "storage.GetTokensSymbols"))
			return
		}

		startTime := time.Now()
		nextRun := startTime.Truncate(5 * time.Minute).Add(5 * time.Minute)
		time.Sleep(time.Until(nextRun))

		updateAt := time.Now()

		var wg sync.WaitGroup
		for _, token := range *tokensTrade {
			wg.Add(1)
			go func(token core.TokenData) {
				defer wg.Done()
				c.collectInfo(token, updateAt)
			}(token)
		}
		wg.Wait()
	}
}

func (c *Checker) getPriceToken(symb string) (float64, error) {
	query := "https://api.mexc.com/api/v3/avgPrice?symbol=" + symb + "USDT" //TODO: подумать над остальными источниками
	resp, err := c.httpClient.Get(query)
	if err != nil {
		c.log.Printf("Ошибка при выполнении HTTP-запроса к API МЭКС: %v", err) //TODO: англ. также подумать на формулировкой
		return 0.0, err
	}
	defer resp.Body.Close()

	if resp.Status[:3] != "200" {
		if resp.Status[:3] == "429" {
			c.log.Printf("429 - Rate limit exceeded") //TODO: англ. также подумать на формулировкой
			c.log.Fatal("Rate limit exceeded")        //TODO: англ. также подумать на формулировкой
			//FIXME: Убрать Fatal
		}
		c.log.Printf(fmt.Sprintf("%s - %s", symb, resp.Status))                      //TODO: англ. также подумать на формулировкой
		c.log.Printf("Неверный статус ответа для символа %s: %s", symb, resp.Status) //TODO: англ. также подумать на формулировкой
	}

	var res core.PriceToken
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		c.log.Printf("Ошибка при декодировании JSON для символа %s: %v", symb, err) //TODO: англ. также подумать на формулировкой
		return 0.0, err
	}

	price, err := strconv.ParseFloat(res.Price, 64)
	if err != nil {
		c.log.Printf("Ошибка при парсинге цены для символа %s: %v", symb, err) //TODO: англ. также подумать на формулировкой
		return 0.0, err
	}

	return price, nil
}

func (c *Checker) collectInfo(token core.TokenData, updatedAt time.Time) {
	newTokenPrice, err := c.getPriceToken(token.Symbol)
	if err != nil {
		c.log.Printf("Ошибка при получении цены токена для символа %s: %v", token.Symbol, err) //TODO: подумать над остальными источниками
		return
	}

	if token.Price == newTokenPrice{
		c.log.Printf("% not change token data %s", token.Price)
		return
	}

	err = c.storage.UpdatePriceToken(token.Symbol)
	if err != nil{
		c.log.Errorf("failed update price token", err) //TODO: англ. также подумать на формулировкой, проверить вывод
	}

	data := &kafkaMessages.Token{
		Symbol: token.Symbol,
		Price: token.Price,
		UpdatedAT: timestamppb.New(updatedAt.UTC()),
	}

	dataBytes, err := proto.Marshal(data)
	if err != nil {
		c.log.Errorf("Ошибка кодирования proto.Marshal", err) //TODO: англ. также подумать на формулировкой, проверить вывод
		return
	}

	msg := kafka.Message{
		Topic: "", //TODO 
		Value: dataBytes,
		Time: time.Now().UTC(),
		// Headers: , //TODO
	}

	c.kafkaClientProducer.PublishMessage(context.Background(), msg)
}
