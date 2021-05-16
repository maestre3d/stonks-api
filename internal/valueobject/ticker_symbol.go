package valueobject

import (
	"strings"

	"github.com/neutrinocorp/ddderr"
)

type TickerSymbol string

const (
	tickerSymbolMinLength = 1
	ticerkSymbolMaxLength = 5
)

var (
	ErrTickerSymbolOutOfRange = ddderr.NewOutOfRange("ticker_symbol", tickerSymbolMinLength, ticerkSymbolMaxLength)
)

func NewTickerSymbol(value string) (TickerSymbol, error) {
	symbol := TickerSymbol(strings.ToUpper(value))
	if err := symbol.ensureValidLength(); err != nil {
		return "", err
	}
	return symbol, nil
}

func (s TickerSymbol) ensureValidLength() error {
	if tickerLength := len(s); tickerLength < tickerSymbolMinLength ||
		tickerLength > ticerkSymbolMaxLength {
		return ErrTickerSymbolOutOfRange
	}
	return nil
}

func (s TickerSymbol) Value() string {
	return string(s)
}
