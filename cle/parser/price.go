package parser

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"tvalert/cle"
	"tvalert/cle/lexer"
)

type PriceType int

const (
	PRICE PriceType = iota
	DIFFERENCE
	PERCENTPRICE
	MARKET
)

type Price struct {
	Type        PriceType
	PriceSource string
	Duration    int64   //Optional
	IsLaddered  [2]bool //Optional
	//0,0 -> no, 1,0 -> laddered; 1,1 -> exponential laddered
	Values [3]float64 // [0] Seperation [1]Value1 [2]Value2
}

func ParsePrice(tk []lexer.Token) (p Price, err error) {
	if len(tk) == 0 {
		return p, nerr(empty, "Error Parse Price no Input")
	}

	p.PriceSource = "close"
	if tk[0].Type == lexer.SOURCE {
		switch tk[0].Content {
		case "high", "low", "close", "open":
			p.PriceSource = tk[0].Content
		default:
			return p, nerr(empty, fmt.Sprintf("Parse Price, Invalid Source Value %s", tk[0].Content))
		}

		if len(tk) < 2 || tk[1].Type != lexer.DURATION {
			return p, nerr(empty, fmt.Sprintf("after -%s a duration has to follow", p.PriceSource))
		}

		ss := tk[1].Content
		n, err := strconv.Atoi(ss[:len(ss)-1])
		if err != nil {
			return p, err
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

	switch tk[0].Type {
	case lexer.FLOAT: // 30000 places order at $30000
		p.Type = PRICE
	case lexer.DFLOAT: // -300 places order $300 below the marketprice
		p.Type = DIFFERENCE
	case lexer.PERCENT: // 2% places order 2% below the marketprice
		p.Type = PERCENTPRICE
	case lexer.MARKET: // -market market buys
		p.Type = MARKET
	case lexer.FLAG: // -l -le for laddered Orders
		err = ParsePriceFlag(&p, tk[0].Content, tk[1:])
		return p, err
	default:
		return p, nerr(empty, fmt.Sprintf("Error Price Parsing, %v %s is not a valid price", tk[0].Type, tk[0].Content))
	}

	p.Values[0], err = strconv.ParseFloat(tk[0].Content, 64)
	if err != nil {
		return p, nerr(err, fmt.Sprintf("Error Price Parsing %s is not a Number", tk[0].Content))
	}

	return p, nil
}

// ParsePriceFlag parses laddered Order
func ParsePriceFlag(p *Price, flag string, tl []lexer.Token) (err error) {
	if len(tl) > 3 {
		return nerr(empty, "Parse Price Flag Error, Not Enough Arguments")
	}

	switch flag {
	case "l": //laddered Order
		p.IsLaddered = [2]bool{true, false}
	case "le": //exponential laddered Order
		p.IsLaddered = [2]bool{true, true}
	default:
		return errors.New("This Flag is not supported: " + flag)
	}

	if len(tl) < 3 {
		return errors.New("Not enough Arguments for a laddered order")
	}

	if tl[0].Type == lexer.FLOAT { //First Value sets up how many orders are placed
		num, err := strconv.Atoi(tl[0].Content)
		if err != nil {
			return err
		}

		if num > 25 || num < 2 {
			return nerr(empty, "Error Parse Price Flag, number of seperation to high, max is 25")
		}
		p.Values[0] = float64(num)
	} else {
		return nerr(empty, fmt.Sprintf("Error Parse Price Flage, First Value %s must be a Number ", tl[0].Content))
	}

	if tl[1].Type != tl[2].Type {
		return nerr(empty, "Values 2 and 3 Arguments must be from same type")
	}

	switch tl[1].Type {
	case lexer.FLOAT:
		p.Type = PRICE
	case lexer.DFLOAT:
		p.Type = DIFFERENCE
	case lexer.PERCENT:
		p.Type = PERCENTPRICE
	default:
		return nerr(empty, fmt.Sprintf("Error Parsing Price Flag! %+v is not a legit Pricevalue", tl[1]))
	}
	v1, err := strconv.ParseFloat(tl[1].Content, 64)
	if err != nil {
		return nerr(err, "Error Parsing Price Flag!")
	}

	if err != nil {
		return err
	}
	v2, err := strconv.ParseFloat(tl[2].Content, 64)

	p.Values[1], p.Values[2] = v1, v2

	return nil
}

func (p *Price) Evaluate(f cle.CLEIO, ws cle.Logger, side string, ticker string, size float64) (err error) {
	var mp float64

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
		err = p.EvaluatePrice(f, side, ticker, size, mp, ws, false)
	case DIFFERENCE:
		err = p.EvaluateDifference(f, side, ticker, size, mp, ws, false)
	case PERCENTPRICE:
		err = p.EvaluatePercentual(f, side, ticker, size, mp, ws, false)
	}

	return err
}

func (p *Price) EvaluatePrice(f cle.CLEIO, sides string, ticker string, size float64, mp float64, ws cle.Logger, ro bool) error {
	var side bool
	if sides == "buy" {
		side = true
	}

	if !p.IsLaddered[0] {
		so, err := f.SetOrder(side, ticker, p.Values[0], size, ro)
		if err != nil {
			return err
		}

		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		ws.Write([]byte(s))

		return err
	}

	p1, p2 := p.Values[1], p.Values[2]
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2) //provides an arry of [2]float64 a pair saves the orders price and position size

	for _, v := range plo {
		so, err := f.SetOrder(side, ticker, v[0], size*v[1], ro)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		ws.Write([]byte(s))
	}
	return nil
}

func (p *Price) EvaluateDifference(f cle.CLEIO, side string, ticker string, size float64, mp float64, ws cle.Logger, ro bool) error {
	factor := 1.0
	ss := true
	if side == "sell" {
		factor = -1.0
		ss = false
	}

	if !p.IsLaddered[0] {
		so, err := f.SetOrder(ss, ticker, mp-p.Values[0]*factor, size, ro)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		ws.Write([]byte(s))
		return err
	}

	p1, p2 := mp-p.Values[1]*factor, mp-p.Values[2]*factor
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2)

	var wg sync.WaitGroup
	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	out := make([][]byte, len(plo), len(plo))

	for i, v := range plo {
		wg.Add(1)
		go func(ii int, vv [2]float64) {
			defer wg.Done()
			so, err := f.SetOrder(ss, ticker, vv[0], size*vv[1], ro)
			if err != nil {
				fatalErrors <- err
			}
			s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
			out[ii] = []byte(s)
		}(i, v)
		/*
			so, err := f.SetOrder(ss, ticker, v[0], size*v[1], ro)
			if err != nil {
				return err
			}
			s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Ticker, so.Size, so.Price)
			ws.Write([]byte(s))
		*/
	}
	//IF
	go func() {
		wg.Wait()
		for _, v := range out {
			ws.Write(v)
		}
		close(wgDone)
	}()
	switch {
	case <-wgDone:

		break
	case nil != <-fatalErrors:
		close(fatalErrors)
		return nerr(empty, "Error Occured while Getting Exchange Data")
	}
	close(fatalErrors)

	return nil
}

func (p *Price) EvaluatePercentual(f cle.CLEIO, side string, ticker string, size float64, mp float64, ws cle.Logger, ro bool) error {
	factor := 1.0
	ss := true
	if side == "sell" {
		factor = -1.0
		ss = false
	}
	//var err error

	if !p.IsLaddered[0] {
		so, err := f.SetOrder(ss, ticker, mp-mp*p.Values[0]/100*factor, size, ro)
		if err != nil {
			return err
		}

		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		ws.Write([]byte(s))
		return err
	}

	p11, p22 := mp*p.Values[1]/100, mp*p.Values[2]/100
	p1, p2 := mp-p11*factor, mp-p22*factor
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2)

	for _, v := range plo {
		so, err := f.SetOrder(ss, ticker, v[0], size*v[1], ro)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("Placed: %s %s %f %f ", so.Side, so.Market, so.Size, so.Price)
		ws.Write([]byte(s))
	}

	return nil
}

func GetPricesLadderedOrder(exponential bool, split, p1, p2 float64) [][2]float64 {
	b := (p2 - p1) / split
	k := b * split / (split - 1)
	//k := (p2 - p1 + b) / split

	sum := (split + 1) / 2
	var fn func(iterate int) float64

	// i(i)
	if exponential {
		fn = func(iterate int) float64 {
			return (float64(iterate+1) / split) / sum
		}
	} else {
		fn = func(iterate int) float64 {
			return 1 / split
		}
	}

	var o [][2]float64
	for i := 0; i < int(split); i++ {
		o = append(o, [2]float64{p1 + k*float64(i), fn(i)})
	}
	return o
}

func harmonicSum(n int) float64 {
	var sum float64
	for i := 0; i < n; i++ {
		sum += 1 / (float64(i) + 1)
	}
	return sum
}
