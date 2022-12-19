package ftx

import (
	"encoding/json"
	"tradingviewListener/exchange"
)

func (p *FTX) OpenOrders() ([]exchange.Order, error) {
	var oor OpenOrdersResponse
	resp, err := p.get("orders", []byte(""))
	if err != nil {
		return oor.Result, err
	}
	err = processResponse(resp, &oor)
	out := oor.Result

	return out, err
}

func (p *FTX) OpenTriggerOrders() ([]exchange.TriggerOrder, error) {
	var oor OpenTriggerOrdersResponse
	resp, err := p.get("conditional_orders", []byte(""))
	if err != nil {
		return oor.Result, err
	}
	err = processResponse(resp, &oor)
	out := oor.Result

	return out, err
}

func (p *FTX) SetOrder(sid bool, ticker string, price float64, size float64, reduceOnly bool) (exchange.Order, error) {
	side := "sell"
	if sid {
		side = "buy"
	}

	var out SetOrderResponse
	rq, err := json.Marshal(NewOrder{
		Market:     ticker,
		Side:       side,
		Price:      price,
		Size:       size,
		ReduceOnly: reduceOnly,
		Type:       "limit",
	})
	if err != nil {
		return out.Result, err
	}
	resp, err := p.post("orders", rq)
	if err != nil {
		return out.Result, err
	}
	err = processResponse(resp, &out)
	if err != nil {
		return out.Result, err
	}
	return out.Result, nil
}

// orderType "stop", "trailingStop", "takeProfit"; default is stop
func (p *FTX) SetTriggerOrder(sid bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (exchange.TriggerOrder, error) {
	side := "sell"
	if sid {
		side = "buy"
	}

	var out SetTriggerOrderResponse
	rq, err := json.Marshal(NewTriggerOrder{
		Market:       ticker,
		Side:         side,
		TriggerPrice: price,
		Size:         size,
		ReduceOnly:   reduceOnly,
		Type:         orderType,
	})
	if err != nil {
		return out.Result, err
	}
	resp, err := p.post("conditional_orders", rq)
	if err != nil {
		return out.Result, err
	}
	err = processResponse(resp, &out)
	if err != nil {
		return out.Result, err
	}
	return out.Result, nil
}
