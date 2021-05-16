package event

type AssetRegistered struct {
	TickerSymbol string  `json:"ticker_symbol"`
	SpreadPrice  float64 `json:"spread_price"`
}

func (a AssetRegistered) AggregateID() string {
	return a.TickerSymbol
}

func (AssetRegistered) Name() string {
	return newAssetEventName("registered")
}
