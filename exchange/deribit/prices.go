package deribit

import (
	"errors"
	"github.com/frankrap/deribit-api/models"
	"strconv"
	"time"
	"tvalert/exchange"
)

func (d *Deribit) Highest(ticker string, duration int64) (float64, error) {
	end := time.Now().UnixMilli()
	start := end - duration*1000
	res := "60"

	if duration > 1440 {
		res = "360"
	}
	if duration > 3000 {
		res = "1D"
	}

	resp, err := d.d.GetTradingviewChartData(&models.GetTradingviewChartDataParams{
		InstrumentName: ticker,
		StartTimestamp: start,
		EndTimestamp:   end,
		Resolution:     res,
	})
	if err != nil {
		return 0, err
	}
	if len(resp.High) == 0 {
		return 0, errors.New("Value Wrong")
	}

	var high float64
	for _, v := range resp.High {
		if v > high {
			high = v
		}
	}
	return high, nil
}

func (d *Deribit) Lowest(ticker string, duration int64) (float64, error) {
	end := time.Now().UnixMilli()
	start := end - duration*1000
	res := "60"

	if duration > 1440 {
		res = "360"
	}
	if duration > 3000 {
		res = "1D"
	}

	resp, err := d.d.GetTradingviewChartData(&models.GetTradingviewChartDataParams{
		InstrumentName: ticker,
		StartTimestamp: start,
		EndTimestamp:   end,
		Resolution:     res,
	})
	if err != nil {
		return 0, err
	}
	if len(resp.Low) == 0 {
		return 0, errors.New("Value Wrong")
	}

	var low float64 = resp.Low[0]
	for _, v := range resp.High {
		if v < low {
			low = v
		}
	}
	return low, nil
}

func (d *Deribit) MarketPrice(ticker string) (float64, error) {
	res, err := d.d.Ticker(&models.TickerParams{InstrumentName: ticker})
	return (res.BestAskPrice + res.BestBidPrice + res.LastPrice) / 3, err
}

func resolutionToString(i int) string {
	switch i {
	case 1440:
		return "D"
	default:
		return strconv.Itoa(i)
	}
}

func checkResolution(res int64) int64 {
	return fnRes(res)
}

var fnRes = exchange.GenerateResolutionFunc(1440, 720, 360, 180, 120, 60, 30, 15, 5, 3, 1)
