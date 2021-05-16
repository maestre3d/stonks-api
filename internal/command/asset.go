package command

import (
	"context"

	"github.com/maestre3d/stonks-api/internal/usecase"
	"github.com/maestre3d/stonks-api/internal/valueobject"
)

type RegisterAsset struct {
	TickerSymbol string  `json:"ticker_symbol"`
	SpreadPrice  float64 `json:"spread_price"`
}

func RegisterAssetHandler(ctx context.Context, u *usecase.RegisterAsset, cmd RegisterAsset) error {
	symbol, err := valueobject.NewTickerSymbol(cmd.TickerSymbol)
	if err != nil {
		return err
	}
	spread, err := valueobject.NewSpread(cmd.SpreadPrice)
	if err != nil {
		return err
	}
	return u.Register(ctx, symbol, spread)
}
