package storage

import "github.com/fishmanDK/price_checker/internal/core"

type Storage interface{
	GetTokensSymbols() (*core.TokensSymbols, error)
	UpdatePriceToken(symb string) error
}