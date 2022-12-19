package exchange

import (
	"fmt"
	"time"
)

// A Ticker represents the latest information of a market pair
type Ticker struct {
	Exchange string `json:"exchange,omitempty"`
	//Name is the Name of the ticker market
	Name string `json:"name,omitempty"`
	//The highest ask price at sell side
	Ask float64 `json:"ask,omitempty"`
	//The highest bid price at buy side
	Bid float64 `json:"bid,omitempty"`
	//The last traded price
	Last float64 `json:"last,omitempty"`
	//Updated
	Updated time.Time `json:"-"`
}

// A Candle represents the OHCLV data of a market
type Candle struct {
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Volume float64 `json:"volume"`
	//Unix Time in Second
	StartTime time.Time `json:"startTime"`
}

type FundingRates struct {
	Ticker    string    `json:"future,omitempty"`
	Rate      float64   `json:"rate,omitempty"`
	StartTime time.Time `json:"time,omitempty"`
}

type Order struct {
	Id         int64     `json:"id,omitempty"`
	Market     string    `json:"market,omitempty"`
	Side       bool      `json:"side,omitempty"` //TODO Fixen
	Size       float64   `json:"size,omitempty"`
	Price      float64   `json:"price,omitempty"`
	ReduceOnly bool      `json:"reduce_only,omitempty"`
	Created    time.Time `json:"created,omitempty"`
	FilledSize float64   `json:"filledSize,omitempty"`
}

type TriggerOrder struct {
	Id         int64     `json:"id,omitempty"`
	Ticker     string    `json:"ticker,omitempty"`
	Side       bool      `json:"side,omitempty"`
	Size       float64   `json:"size,omitempty"`
	Price      float64   `json:"price,omitempty"`
	ReduceOnly bool      `json:"reduce_only,omitempty"`
	Created    time.Time `json:"created,omitempty"`
}

type Position struct {
	Side             bool    `json:"side"`
	Future           string  `json:"future"`
	NotionalSize     float64 `json:"cost"`
	PositionSize     float64 `json:"size"`
	UPNL             float64 `json:"unrealizedPnl"`
	PNL              float64 `json:"realizedPnl"`
	EntryPrice       float64 `json:"entryPrice"`
	LiquidationPrice float64 `json:"estimatedLiquidationPrice"`
	AvgOpen          float64 `json:"recentAverageOpenPrice"`
	BreakEven        float64 `json:"recentBreakEvenPrice"`
}

type Coin struct {
	Coin     string  `json:"coin,omitempty"`
	Free     float64 `json:"free,omitempty"`
	Total    float64 `json:"total,omitempty"`
	UsdValue float64 `json:"usdValue,omitempty"`
}

type Fill struct {
	Fee          float64   `json:"fee,omitempty"`
	FeeCurrency  string    `json:"feeCurrency,omitempty"`
	Future       string    `json:"market,omitempty"`
	Id           int64     `json:"id,omitempty"`
	OrderId      int       `json:"orderId,omitempty"`
	Price        float64   `json:"price,omitempty"`
	BaseCurrency string    `json:"baseCurrency,omitempty"`
	Side         string    `json:"side,omitempty"`
	Size         float64   `json:"size,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}

type Account struct {
	Username          string     `json:"username"`
	Collateral        float64    `json:"collateral"`
	FreeCollateral    float64    `json:"freeCollateral"`
	Leverage          float64    `json:"leverage"`
	MarginFraction    float64    `json:"marginFraction"`
	TotalAccountValue float64    `json:"totalAccountValue"`
	TotalPositionSize float64    `json:"totalPositionSize"`
	Positions         []Position `json:"positions"`
}

// New resolution must me greater than old
func ConvertChartResolution(startResolution, GoalResolution int64, ch []Candle) ([]Candle, error) {
	if GoalResolution == startResolution {
		return ch, nil
	}

	if startResolution > GoalResolution || GoalResolution%startResolution != 0 {
		return ch, fmt.Errorf("New Res %v and old %v do not fit", GoalResolution, startResolution)
	}

	quotient := int(GoalResolution / startResolution)

	var newChart []Candle = make([]Candle, 0, len(ch)/quotient)

	for _, c := range ch {
		if c.StartTime.Unix()%GoalResolution != 0 {
			ch = ch[1:]
		} else {
			break
		}
	}

	for {
		if len(ch) < quotient {
			break
		}
		newChart = append(newChart, ConvertCandleResolution(ch[:quotient]))
		ch = ch[quotient:]
	}
	if len(ch) != 0 {
		newChart = append(newChart, ConvertCandleResolution(ch))
	}

	return newChart, nil
}

// ConvertResolution converts the a lower resolution into a higher resolution
func ConvertCandleResolution(c []Candle) Candle {
	var out Candle = Candle{c[0].Close, c[0].High, c[0].Low, c[0].Open, c[0].Volume, c[0].StartTime}

	if len(c) == 1 {
		return c[0]
	}

	for _, i := range c[1:] {
		out.Close = i.Close
		out.Volume += i.Volume
		if i.High > out.High {
			out.High = i.High
		}
		if i.Low < out.Low {
			out.Low = i.Low
		}
	}
	return out
}

func GenerateResolutionFunc(resInSeconds ...int64) func(int64) int64 {
	return func(r int64) int64 {
		var newRes int64
		for _, v := range resInSeconds {
			if r == v {
				newRes = r
				return newRes
			}
		}
		for _, v := range resInSeconds {
			if r >= v && r%v == 0 {
				return v
			}
		}
		return 3600
	}
}
