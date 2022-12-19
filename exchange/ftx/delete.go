package ftx

import (
	"encoding/json"
	"fmt"
)

type RequestDeleteOrder struct {
	Market                string `json:"market,omitempty"`
	Side                  string `json:"side,omitempty"`
	LimitOnly             bool   `json:"limitOnly,omitempty"`
	ConditionalOrdersOnly bool   `json:"conditionalOrdersOnly,omitempty"`
}

func (p *FTX) Cancel(Side string, Ticker string) error {
	rB, err := json.Marshal(RequestDeleteOrder{
		Market:                Ticker,
		Side:                  Side,
		LimitOnly:             true,
		ConditionalOrdersOnly: false,
	})
	if err != nil {
		return err
	}
	_, err = p.delete("orders", rB)
	return err
}

func (p *FTX) CancelTrigger(Side string, Ticker string) error {
	rB, err := json.Marshal(RequestDeleteOrder{
		Market:                Ticker,
		Side:                  Side,
		ConditionalOrdersOnly: true,
	})
	if err != nil {
		return err
	}

	_, err = p.delete("orders", rB)
	return err
}

func (p *FTX) CancelOrderById(id int64) error {
	_, err := p.delete(fmt.Sprintf("orders/%d", id), nil)
	return err
}

func (p *FTX) CancelTriggerOrderById(id int64) error {
	_, err := p.delete(fmt.Sprintf("conditional_orders/%d", id), nil)
	return err
}
