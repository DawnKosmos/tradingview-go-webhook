package ftx

import (
	"errors"
	"strconv"
	"tradingviewListener/exchange"
)

func (p *FTX) Market(ticker string) (exchange.Ticker, error) {
	var r MarketResponseSingle
	resp, err := p.get("markets/"+ticker, []byte(""))
	if err != nil {
		return exchange.Ticker{}, err
	}
	err = processResponse(resp, &r)
	if err != nil {
		return exchange.Ticker{}, errors.New("Ticker does not exist")
	}
	t := r.Result

	return exchange.Ticker{
		Exchange: "ftx",
		Name:     t.Name,
		Ask:      t.Ask,
		Bid:      t.Bid,
		Last:     t.Last,
	}, err
}

func (p *FTX) Candles(ticker string, resolution int64, st int64, et int64) ([]exchange.Candle, error) {
	var candles []exchange.Candle
	var end int64 = 0
	newResolution := checkResolution(resolution)

	for st < et {
		end = st + 1500*newResolution
		if end >= et {
			c, err := p.historicalPrices(ticker, newResolution, st, et)
			if err != nil {
				return candles, err
			}
			candles = append(candles, c...)
		} else {
			c, err := p.historicalPrices(ticker, newResolution, st, end)
			if err != nil {
				return candles, err
			}
			candles = append(candles, c...)
		}
		st = st + 1501*newResolution
	}

	formatted, err := exchange.ConvertChartResolution(newResolution, resolution, candles)

	return formatted, err
}

// checkResolution looking if the asked resolution is a valid one
func checkResolution(res int64) int64 {
	var newRes int64
	if res == 3600 || res == 14400 || res == 86400 || res == 300 || res == 60 || res == 900 {
		newRes = res
		return newRes
	}
	if res >= 86400 && res%86400 == 0 {
		return 86400
	}
	if res >= 14400 && res%14400 == 0 {
		return 14400
	}
	if res >= 3600 && res%3600 == 0 {
		return 3600
	}
	if res >= 900 && res%900 == 0 {
		return 900
	}

	if res >= 300 && res%300 == 0 {
		return 300
	}

	if res >= 15 && res%15 == 0 {
		return 15
	}
	return 3600
}

func (p *FTX) historicalPrices(ticker string, resolution int64, st int64, et int64) ([]exchange.Candle, error) {
	var cr CandlesResponse
	resp, err := p.get(
		"markets/"+ticker+
			"/candles?resolution="+strconv.FormatInt(resolution, 10)+
			"&start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		return cr.Result, err
	}
	err = processResponse(resp, &cr)
	return cr.Result, nil
}
