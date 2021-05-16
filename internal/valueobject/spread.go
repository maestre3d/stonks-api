package valueobject

import (
	"github.com/neutrinocorp/ddderr"
)

type Spread float64

const (
	spreadMinRange = 0
)

var ErrSpreadOutOfRange = ddderr.NewOutOfRange("spread", 0, 10000000000000)

func NewSpread(value float64) (Spread, error) {
	spread := Spread(value)
	if err := spread.ensureValidRange(); err != nil {
		return 0, err
	}
	return spread, nil
}

func (s Spread) ensureValidRange() error {
	if s < spreadMinRange {
		return ErrSpreadOutOfRange
	}
	return nil
}

func (s Spread) Value() float64 {
	return float64(s)
}
