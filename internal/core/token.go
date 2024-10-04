package core

type (
	TokensSymbols []TokenData

	TokenData struct {
		Price  float64
		Symbol string
	}

	PriceToken struct {
		Price string `json:"price"`
	}
)
