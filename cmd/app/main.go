package main

import (
	"log"

	"github.com/fishmanDK/price_checker/internal/app"
	"github.com/fishmanDK/price_checker/internal/config"
	"github.com/fishmanDK/price_checker/internal/logger"
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

func main() {
	cfg, err := config.InitConfig()
	if err != nil{
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	app := app.NewApp(appLogger, cfg)
	app.Run()
}