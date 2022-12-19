package cle

import (
	"io"
	"tvalert/cle/lexer"
	"tvalert/cle/parser"
	"tvalert/exchange"
)

type Logger interface {
	io.Writer
	ErrorMessage(s error)
}

func Execute(ex Logger, exchange CLEIO, commands []string) {
	for _, v := range commands {
		tl, err := lexer.Lexer(v)
		if err != nil {
			ex.ErrorMessage(err)
			return
		}
		parsed, err := parser.Parse(tl, ex)
		if err != nil {
			ex.ErrorMessage(err)
			return
		}
		err = parsed.Evaluate(ex, exchange)
		if err != nil {
			ex.ErrorMessage(err)
			return
		}
	}

	err := exchange.ResetTempVariables()
	if err != nil {
		ex.ErrorMessage(err)
	}
}

type CLEIO interface {
	//Name returns the Exchange Name and Subaccount Name if available
	Name() string
	//SetOrder sets an Order, reduceOnly optional, default is false. Returns the set order
	SetOrder(side bool, ticker string, price float64, size float64, reduceOnly bool) (exchange.Order, error)
	//OpenOrders Returns open orders for given ticker
	OpenOrders(side bool, ticker string) ([]exchange.Order, error)
	//SetTriggerOrder set an TriggerOrder, reduceOnly optional, default is true
	SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (exchange.TriggerOrder, error)
	//MarketPrice return the Market Price of the asked Ticker
	MarketPrice(ticker string) (float64, error)
	//Cancel All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	Cancel(Side int64, Ticker string) error
	//CancelTrigger All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	CancelTrigger(Side int64, Ticker string) error
	//Highest return the Highest Price of the ticker for the given duration
	Highest(ticker string, duration int64) (float64, error)
	//Lowest returns the Lowest Price of the ticker for the given duration
	Lowest(ticker string, duration int64) (float64, error)
	//FreeCollateral returns the amount of free collateral in USD
	FreeCollateral(ticker string) (float64, error)
	//OpenPositions returns all Open positions
	OpenPositions() (map[string]exchange.Position, error)
	ResetTempVariables() error
}
