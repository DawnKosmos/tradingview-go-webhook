package deribit

import (
	"github.com/frankrap/deribit-api/models"
	"strings"
	"tvalert/exchange"
)

func (d *Deribit) FreeCollateral(ticker string) (float64, error) {
	var enum string
	var price float64
	var err error
	switch {
	case strings.Contains(ticker, "USDC"):
		enum = "USDC"
		price = 1
	case strings.Contains(ticker, "BTC"):
		enum = "BTC"
		price, err = d.indexPrice("btc_usd")
		if err != nil {
			return 0, err
		}
	case strings.Contains(ticker, "ETH"):
		enum = "ETH"
		price, err = d.indexPrice("eth_usd")
		if err != nil {
			return 0, err
		}
	case strings.Contains(ticker, "SOL"):
		enum = "SOL"
		price, err = d.indexPrice("sol_usd")
		if err != nil {
			return 0, err
		}
	default:
		enum = "USDC"
	}

	res, err := d.d.GetAccountSummary(&models.GetAccountSummaryParams{
		Currency: enum,
		Extended: false,
	})

	return res.AvailableFunds * price, err
}

func (d *Deribit) OpenPositions() (map[string]exchange.Position, error) {
	var currencies []string = []string{"BTC", "ETH", "SOL", "USDC"}
	position := make(map[string]exchange.Position)
	for _, v := range currencies {
		pos, err := d.d.GetPositions(&models.GetPositionsParams{
			Currency: v,
		})
		if err != nil {
			return nil, err
		}
		for _, vv := range pos {
			var side bool
			if vv.Direction == "buy" {
				side = true
			}
			position[vv.InstrumentName] = exchange.Position{
				Side:             side,
				Future:           vv.InstrumentName,
				NotionalSize:     vv.Size,
				PositionSize:     vv.SizeCurrency,
				UPNL:             vv.TotalProfitLoss,
				PNL:              vv.RealizedProfitLoss,
				EntryPrice:       vv.AveragePrice,
				LiquidationPrice: vv.EstimatedLiquidationPrice,
				AvgOpen:          vv.AveragePrice,
				BreakEven:        vv.AveragePrice, //TODO Calculation
			}

		}
	}
	return position, nil
}

func (d *Deribit) ResetTempVariables() error {
	return nil
}

type indexRequest struct {
	IndexName string `json:"index_name"`
}

type indexResponse struct {
	IndexPrice float64 `json:"index_price"`
}

func (d *Deribit) indexPrice(s string) (float64, error) {
	var out indexResponse
	err := d.d.Call("public/get_index_price", indexRequest{s}, &out)
	return out.IndexPrice, err
}
