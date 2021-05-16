package repository

import (
	"context"

	"github.com/maestre3d/stonks-api/internal/aggregate"
	"github.com/maestre3d/stonks-api/internal/valueobject"
)

type Asset interface {
	Save(context.Context, aggregate.Asset) error
	Update(context.Context, aggregate.Asset) error
	Remove(context.Context, aggregate.Asset) error
	Find(context.Context, valueobject.TickerSymbol) (*aggregate.Asset, error)
	Search(context.Context) ([]*aggregate.Asset, error)
}
