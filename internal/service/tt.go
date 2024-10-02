package service

import (
	"os"
	"time"
)

var errorsFile, _ = os.Open("data/errors.txt")

const (
	keyInfluxDB = "VRsOE9gM1OcQS5VFEwiR3pZERheEi4fXm909k_C7aaALGql509v3nNH4oWe5ziW70K0bXDG-P3Js-YxJfndHFA=="
	org = "price_checker"
	bucket = "tokens"
)
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

type SymbolPriceToken struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

var tokensTrade = map[string]struct{}{
	"BTC":        struct{}{},
	"ETH":        struct{}{},
	"ADA":        struct{}{},
	"BNB":        struct{}{},
	"SOL":        struct{}{},
	"USDC":       struct{}{},
	"XRP":        struct{}{},
	"DOGE":       struct{}{},
	"TON":        struct{}{},
	"TRX":        struct{}{},
	"AVAX":       struct{}{},
	"SHIB":       struct{}{},
	"LINK":       struct{}{},
	"BCH":        struct{}{},
	"DOT":        struct{}{},
	"DAI":        struct{}{},
	"LEO":        struct{}{},
	"LTC":        struct{}{},
	"NEAR":       struct{}{},
	"SUI":        struct{}{},
	"APE":        struct{}{},
	"UNI":        struct{}{},
	"CRV":        struct{}{},
	"ATOM":       struct{}{},
	"HBAR":       struct{}{},
	"MANA":       struct{}{},
	"FTM":        struct{}{},
	"RUNE":       struct{}{},
	"GLMR":       struct{}{},
	"ICP":        struct{}{},
	"RAY":        struct{}{},
	"CAKE":       struct{}{},
	"KSM":        struct{}{},
	"FIL":        struct{}{},
	"AXS":        struct{}{},
	"YFI":        struct{}{},
	"DYDX":       struct{}{},
	"OP":         struct{}{},
	"KLAY":       struct{}{},
	"CRO":        struct{}{},
	"ZEC":        struct{}{},
	"COMP":       struct{}{},
	"SNX":        struct{}{},
	"AAVE":       struct{}{},
	"MKR":        struct{}{},
	"XTZ":        struct{}{},
	"ALGO":       struct{}{},
	"EOS":        struct{}{},
	"GRT":        struct{}{},
	"NEO":        struct{}{},
	"SC":         struct{}{},
	"ONT":        struct{}{},
	"QTUM":       struct{}{},
	"ETC":        struct{}{},
	"KAVA":       struct{}{},
    "WAVES":      struct{}{},
    "CELO":       struct{}{},
    "HNT":        struct{}{},
    "NKN":        struct{}{},
    "ZIL":        struct{}{},
    "WAXP":       struct{}{},
    "CHZ":        struct{}{},
	"AIC":        struct{}{},
	"SEI":        struct{}{},
    "PEPE":       struct{}{},
    "TOMI":       struct{}{},
	"FET":        struct{}{},
	"NEIRO":        struct{}{},
	"SHIT":        struct{}{},
	"CATI":        struct{}{},
	"DOGS":        struct{}{},
	"FIDA":        struct{}{},
	"ZK":        struct{}{},
	"JASMY":        struct{}{},
	"KAS":        struct{}{},
	"RVN":        struct{}{},
	"WOO":        struct{}{},
	"APT":        struct{}{},
	"PEOPLE":        struct{}{},
	"HOT":        struct{}{},
	"1INCH":        struct{}{},
	"NEXO":        struct{}{},
	"RPL":        struct{}{},
}

type StatisticAlong struct {
	// Left int
	// Right int
	Current  int
	TotalSum float64
}

type StatisticToken struct {
	PricesAlong1H    StatisticAlong
	PricesAlont30Min StatisticAlong

	Round1H     int
	Prices1H    [12]float64
	Round30Min  int
	Prices30Min [6]float64

	Start       time.Time
	PriceStart  float64
	Coefficient float64
	LastPrice   float64
}