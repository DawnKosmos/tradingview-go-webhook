package parser

import (
	"fmt"
	"strconv"
	"tvalert/cle"
	"tvalert/cle/lexer"
)

type StopOrder struct {
	Side   string
	Ticker string
	A      Amount
	P      StopPrice
}

func ParseStop(tk []lexer.Token) (o *StopOrder, err error) {
	o = &StopOrder{}
	if len(tk) < 4 {
		return nil, nerr(empty, "Error Parse Stop, Input not enough arguments")
	}

	var a Amount

	if tk[0].Type != lexer.SIDE {
		return nil, nerr(empty, "Error Parse Stop, After Stop buy/sell has to be provided")
	}
	o.Side = tk[0].Content

	if tk[1].Type == lexer.VARIABLE {
		o.Ticker = tk[1].Content
	} else {
		return nil, nerr(empty, fmt.Sprintf("Error Parse Stop, %s is no Ticker", tk[0].Content))
	}

	switch tk[2].Type {
	case lexer.FLOAT: // 5.2 -> 5.2 Coins
		a.Type = COIN
	case lexer.UFLOAT: // u500 -> 500 USD worth of the coin
		a.Type = FIAT
	case lexer.PERCENT: // 100% -> 100% of your Free Collateral of the Coin
		a.Type = ACCOUNTSIZE
	case lexer.POSITION: // -position -> 100% of the Positions Size
		a.Type = POSITIONSIZE
	case lexer.PRICEORDER:
		a.Type = ALL
	default:
		return nil, nerr(empty, fmt.Sprintf("Error Parse Stop Order, false Order Size of type"))
	}

	a.Ticker = o.Ticker

	a.Value, err = strconv.ParseFloat(tk[2].Content, 64)
	if err != nil {
		return nil, nerr(err, fmt.Sprintf("Erorr Parse Stop Wrong Value should be a Float is %s", tk[2].Content))
	}

	o.A = a

	o.P, err = ParseStopPrice(tk[3:])
	return o, err
}

// STOPPRICE

type StopPrice struct {
	Type        PriceType
	PriceSource string
	Duration    int64
	Value       float64
}

func ParseStopPrice(tk []lexer.Token) (p StopPrice, err error) {
	if len(tk) == 0 {
		return p, nerr(empty, "Error Parse Stop Price: no Input")
	}
	p.PriceSource = "close"
	if tk[0].Type == lexer.SOURCE {
		switch tk[0].Content {
		case "high", "low", "close", "open":
			p.PriceSource = tk[0].Content
		default:
			return p, nerr(empty, fmt.Sprintf("Parse Stop Price, Invalid Source Value %s", tk[0].Content))
		}

		if len(tk) < 2 || tk[1].Type != lexer.DURATION {
			return p, nerr(empty, fmt.Sprintf("Error Parse Stop Price: after -%s a duration has to follow", p.PriceSource))
		}

		ss := tk[1].Content
		n, err := strconv.Atoi(ss[:len(ss)-1])
		if err != nil {
			return p, nerr(err, "Error Parse Stop Price: Invalid Duration")
		}

		switch ss[len(ss)-1] {
		case 'h':
			n *= 3600
		case 'm':
			n *= 60
		case 'd':
			n *= 3600 * 24
		default:
			return p, nerr(empty, fmt.Sprintf("Error Price Parsing Duration with %s !!", ss))
		}
		p.Duration = int64(n)

		tk = tk[2:]
	}

	if len(tk) == 0 {
		return p, nerr(empty, "Error Parsing Price no Price provided")
	}

	switch tk[0].Type {
	case lexer.FLOAT: // 30000 places order at $30000
		p.Type = PRICE
	case lexer.DFLOAT: // -300 places order $300 below the marketprice
		p.Type = DIFFERENCE
	case lexer.PERCENT: // 2% places order 2% below the marketprice
		p.Type = PERCENTPRICE
	default:
		return p, nerr(empty, fmt.Sprintf("Error Stop Price Parsing, %v %s is not a valid price", tk[0].Type, tk[0].Content))
	}

	p.Value, err = strconv.ParseFloat(tk[0].Content, 64)

	if err != nil {
		return p, nerr(err, fmt.Sprintf("Error Price Parsing %s is not a Number", tk[0].Content))
	}

	return p, nil
}

func (o *StopOrder) Evaluate(w cle.Logger, f cle.CLEIO) error {
	size, err := o.A.EvaluateStop(o.Side, f)
	if err != nil {
		return err
	}

	return o.P.Evaluate(w, f, o.Side, o.Ticker, size)
}

func (p *StopPrice) Evaluate(ws cle.Logger, f cle.CLEIO, sides string, ticker string, size float64) (err error) {
	var side bool
	factor := 1.0
	if sides == "buy" {
		side = true
		factor = -1.0
	}

	var mp float64 //marketprice
	switch p.PriceSource {
	case "market", "open", "close":
		mp, err = f.MarketPrice(ticker)
	case "low":
		mp, err = f.Lowest(ticker, p.Duration)
	case "high":
		mp, err = f.Highest(ticker, p.Duration)
	}
	if err != nil {
		return err
	}

	switch p.Type {
	case PRICE:
		so, err := f.SetTriggerOrder(side, ticker, p.Value, size, "stop", true)
		if err != nil {
			return err
		}

		s := fmt.Sprintf("Placed Stop: %s %s %f %f ", sides, so.Ticker, so.Size, so.Price)
		ws.Write([]byte(s))
		return err
	case DIFFERENCE:
		so, err := f.SetTriggerOrder(side, ticker, mp-p.Value*factor, size, "stop", true)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed Stop: %s %s %f %f ", sides, so.Ticker, so.Size, so.Price)
		ws.Write([]byte(s))
		return err

	case PERCENTPRICE:
		so, err := f.SetTriggerOrder(side, ticker, mp-mp*p.Value/100*factor, size, "stop", true)
		if err != nil {
			return err
		}

		s := fmt.Sprintf("Placed: %s %s %f %f ", sides, so.Ticker, so.Size, so.Price)
		ws.Write([]byte(s))
		return err
	}
	return nil
}
