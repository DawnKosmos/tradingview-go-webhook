package ftx

import "tradingviewListener/exchange"

type OpenOrdersResponse struct {
	Success bool             `json:"success"`
	Result  []exchange.Order `json:"result"`
}

type OpenTriggerOrdersResponse struct {
	Success bool                    `json:"success"`
	Result  []exchange.TriggerOrder `json:"result"`
}

type SetOrderResponse struct {
	Success bool           `json:"success"`
	Result  exchange.Order `json:"result"`
}

type PositionsResponse struct {
	Success bool                `json:"success"`
	Result  []exchange.Position `json:"result"`
}

type TriggerOrdersResponse struct {
	Success bool                    `json:"success"`
	Result  []exchange.TriggerOrder `json:"result"`
}

type SetTriggerOrderResponse struct {
	Success bool                  `json:"success"`
	Result  exchange.TriggerOrder `json:"result"`
}

type AccountResponse struct {
	Success bool             `json:"success"`
	Result  exchange.Account `json:"result"`
}

type FillResponse struct {
	Success bool            `json:"success"`
	Result  []exchange.Fill `json:"result"`
}

type WalletResponse struct {
	Success bool            `json:"success"`
	Result  []exchange.Coin `json:"result"`
}

type NewOrder struct {
	Market     string  `json:"market"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Type       string  `json:"type"`
	Size       float64 `json:"size"`
	ReduceOnly bool    `json:"reduceOnly"`
	Ioc        bool    `json:"ioc"`
	PostOnly   bool    `json:"postOnly"`
}

type NewTriggerOrder struct {
	Market       string  `json:"market"`
	Side         string  `json:"side"`
	Size         float64 `json:"size"`
	Type         string  `json:"type"`
	ReduceOnly   bool    `json:"reduceOnly"`
	TriggerPrice float64 `json:"triggerPrice"`
}

type MarketResponse struct {
	Success bool     `json:"success,omitempty"`
	Result  []Ticker `json:"result,omitempty"`
}

type MarketResponseSingle struct {
	Success bool   `json:"success,omitempty"`
	Result  Ticker `json:"result,omitempty"`
}

type Ticker struct {
	//Name is the Name of the ticker market
	Name string `json:"name,omitempty"`
	//The highest ask price at sell side
	Ask float64 `json:"ask,omitempty"`
	//The highest bid price at buy side
	Bid float64 `json:"bid,omitempty"`
	//The last traded price
	Last float64 `json:"last,omitempty"`
}

type CandlesResponse struct {
	Success bool              `json:"success"`
	Result  []exchange.Candle `json:"result"`
}
