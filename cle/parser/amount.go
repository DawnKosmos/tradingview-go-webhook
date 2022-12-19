package parser

import (
	"math"
	"strings"
	"tvalert/cle"
)

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNTSIZE
	POSITIONSIZE
	ALL
)

type Amount struct {
	Ticker string
	Type   AmountType
	Value  float64
}

func (a *Amount) Evaluate(f cle.CLEIO) (float64, error) {
	switch a.Type {
	case COIN:
		return a.Value, nil
	case FIAT:
		m, err := f.MarketPrice(a.Ticker)
		return a.Value / m, err
	case ACCOUNTSIZE:
		m, err := f.MarketPrice(a.Ticker)
		if err != nil {
			return a.Value, err
		}
		collateral, err := f.FreeCollateral(a.Ticker)
		az := collateral * a.Value / 100
		return az / m, nil
	case POSITIONSIZE:
		pz, err := f.OpenPositions()
		tick := strings.ToUpper(a.Ticker)
		positionSize, ok := pz[tick]
		if !ok {
			return 0, nil
		}
		return positionSize.PositionSize, err
	}
	return 0, nerr(empty, "Error Evaluating Amount of Order")
}

func (a *Amount) EvaluateStop(sidestring string, f cle.CLEIO) (float64, error) {
	var side bool
	if sidestring == "buy" {
		side = true
	}

	switch a.Type {
	case COIN:
		return a.Value, nil
	case ACCOUNTSIZE:
		m, err := f.MarketPrice(a.Ticker)
		if err != nil {
			return a.Value, err
		}
		collateral, err := f.FreeCollateral(a.Ticker)
		az := collateral * a.Value / 100
		return az / m, nil
	case POSITIONSIZE:
		pz, err := f.OpenPositions()
		tick := strings.ToUpper(a.Ticker)
		positionSize, ok := pz[tick]
		if !ok {
			return 0, nil
		}
		return positionSize.PositionSize, err
	case ALL:
		var size float64
		pz, err := f.OpenPositions()
		if err != nil {
			return 0, err
		}

		tick := strings.ToUpper(a.Ticker)
		ps, ok := pz[tick]
		if ok {
			size = ps.PositionSize
		}

		orders, err := f.OpenOrders(!side, a.Ticker)
		if err != nil {
			return 0, err
		}

		for _, v := range orders {
			size += math.Abs(v.Size)
		}
		return size, err
	}
	return 0, nerr(empty, "Error Evaluating Amount of Order")
}
