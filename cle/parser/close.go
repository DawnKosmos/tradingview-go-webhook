package parser

import (
	"fmt"
	"strings"
	"tvalert/cle"
	"tvalert/cle/lexer"
)

type Close struct {
	Ticker string
	Side   string
}

func ParseClose(tk []lexer.Token) (o *Close, err error) {
	o = &Close{}
	if len(tk) < 2 {
		return o, nerr(empty, "Error Close Arguments missing")
	}
	if tk[0].Type == lexer.SIDE {
		o.Side = tk[0].Content
	}

	o.Ticker = tk[1].Content
	return
}

func (o *Close) Evaluate(w cle.Logger, f cle.CLEIO) error {
	pz, err := f.OpenPositions()
	if err != nil {
		return err
	}

	tick := strings.ToUpper(o.Ticker)
	pos, ok := pz[tick]
	if !ok {
		w.Write([]byte("no position to close"))
		return nil
	}

	if pos.Side != o.Side {
		return nil
	}

	price := pos.NotionalSize / pos.PositionSize
	if pos.Side == "buy" {
		price -= 0.1 * price
		so, err := f.SetOrder(false, o.Ticker, price, pos.PositionSize, true)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		w.Write([]byte(s))
		return err
	} else {

		price *= -1
		price += 0.1 * price
		so, err := f.SetOrder(true, o.Ticker, price, pos.PositionSize, true)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		w.Write([]byte(s))
		return err
	}

}
