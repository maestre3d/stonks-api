package aggregate

import (
	"github.com/maestre3d/stonks-api/internal/event"
	"github.com/maestre3d/stonks-api/internal/valueobject"
)

// Asset is a financial instrument from a stock market (e.g. Holder companies, Exchange-Traded Funds, Bonds)
type Asset struct {
	TickerSymbol valueobject.TickerSymbol // high-cardinality already, custom ID not required
	SpreadPrice  valueobject.Spread

	events []event.DomainEvent
}

func NewAsset(ticker valueobject.TickerSymbol, spread valueobject.Spread) *Asset {
	asset := &Asset{
		TickerSymbol: ticker,
		SpreadPrice:  spread,
		events:       make([]event.DomainEvent, 0),
	}

	asset.recordEvents(event.AssetRegistered{
		TickerSymbol: asset.TickerSymbol.Value(),
		SpreadPrice:  asset.SpreadPrice.Value(),
	})
	return asset
}

func (a *Asset) recordEvents(events ...event.DomainEvent) {
	a.events = append(a.events, events...)
}

func (a *Asset) PullEvents() []event.DomainEvent {
	pulledEvents := a.events
	a.events = make([]event.DomainEvent, 0)
	return pulledEvents
}
