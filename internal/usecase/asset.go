package usecase

import (
	"context"

	"github.com/maestre3d/stonks-api/internal/aggregate"
	"github.com/maestre3d/stonks-api/internal/event"
	"github.com/maestre3d/stonks-api/internal/repository"
	"github.com/maestre3d/stonks-api/internal/valueobject"
	"github.com/neutrinocorp/ddderr"
)

type RegisterAsset struct {
	assetRepository repository.Asset
	bus             event.Bus
}

func NewRegisterAsset(r repository.Asset, b event.Bus) *RegisterAsset {
	return &RegisterAsset{
		assetRepository: r,
		bus:             b,
	}
}

var (
	ErrAssetAlreadyExists = ddderr.NewAlreadyExists("asset")
)

func (r RegisterAsset) Register(ctx context.Context, ticker valueobject.TickerSymbol, spread valueobject.Spread) error {
	if asset, _ := r.assetRepository.Find(ctx, ticker); asset != nil {
		return ErrAssetAlreadyExists
	}

	asset := aggregate.NewAsset(ticker, spread)
	if err := r.assetRepository.Save(ctx, *asset); err != nil {
		return err
	}
	return r.bus.Publish(ctx, asset.PullEvents()...)
}
