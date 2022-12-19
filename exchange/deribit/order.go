package deribit

import (
	"github.com/frankrap/deribit-api/models"
	"strconv"
	"time"
	"tvalert/exchange"
)

const TRIGGER = "trigger"

func (d *Deribit) SetOrder(side bool, ticker string, price float64, size float64, reduceOnly bool) (exchange.Order, error) {
	if side {
		resp, err := d.d.Buy(&models.BuyParams{
			InstrumentName: ticker,
			Amount:         size,
			Label:          ticker + "buy",
			Price:          price,
			ReduceOnly:     reduceOnly,
		})
		if err != nil {
			return exchange.Order{}, err
		}
		o := resp.Order

		id, _ := strconv.ParseInt(o.OrderID, 10, 64)

		return exchange.Order{
			Id:         id,
			Market:     o.InstrumentName,
			Side:       true,
			Size:       o.Amount,
			Price:      o.Price.ToFloat64(),
			ReduceOnly: o.ReduceOnly,
			Created:    time.Unix(o.CreationTimestamp/1000, 0),
			FilledSize: o.FilledAmount,
		}, nil
	} else {
		resp, err := d.d.Sell(&models.SellParams{
			InstrumentName: ticker,
			Amount:         size,
			Label:          ticker + "sell",
			Price:          price,
			ReduceOnly:     reduceOnly,
		})
		if err != nil {
			return exchange.Order{}, err
		}
		o := resp.Order

		id, _ := strconv.ParseInt(o.OrderID, 10, 64)

		return exchange.Order{
			Id:         id,
			Market:     o.InstrumentName,
			Side:       false,
			Size:       o.Amount,
			Price:      o.Price.ToFloat64(),
			ReduceOnly: o.ReduceOnly,
			Created:    time.Unix(o.CreationTimestamp/1000, 0),
			FilledSize: o.FilledAmount,
		}, nil
	}
}

func (d *Deribit) OpenOrders(side bool, ticker string) ([]exchange.Order, error) {
	var orders []exchange.Order
	o, err := d.d.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{InstrumentName: ticker, Type: "limit"})
	if err != nil {
		return orders, err
	}

	for _, v := range o {
		var ss bool = v.Direction == "buy"
		id, _ := strconv.ParseInt(v.OrderID, 10, 64)

		if ss && side {
			orders = append(orders, exchange.Order{
				Id:         id,
				Market:     v.InstrumentName,
				Side:       side,
				Size:       v.Amount,
				Price:      v.Price.ToFloat64(),
				ReduceOnly: v.ReduceOnly,
				Created:    time.Unix(v.CreationTimestamp/1000, 0),
				FilledSize: v.FilledAmount,
			})
		}
		if !ss && !side {
			orders = append(orders, exchange.Order{
				Id:         id,
				Market:     v.InstrumentName,
				Side:       side,
				Size:       v.Amount,
				Price:      v.Price.ToFloat64(),
				ReduceOnly: v.ReduceOnly,
				Created:    time.Unix(v.CreationTimestamp/1000, 0),
				FilledSize: v.FilledAmount,
			})
		}
	}
	return orders, nil
}

func (d *Deribit) SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (exchange.TriggerOrder, error) {
	if side {
		resp, err := d.d.Buy(&models.BuyParams{
			InstrumentName: ticker,
			Amount:         size,
			Label:          ticker + TRIGGER + "buy",
			Price:          price,
			ReduceOnly:     reduceOnly,
			Type:           "stop_market",
		})
		if err != nil {
			return exchange.TriggerOrder{}, err
		}
		o := resp.Order

		id, _ := strconv.ParseInt(o.OrderID, 10, 64)

		return exchange.TriggerOrder{
			Id:         id,
			Ticker:     ticker,
			Side:       side,
			Size:       size,
			Price:      price,
			ReduceOnly: reduceOnly,
			Created:    time.Unix(o.CreationTimestamp, 0),
		}, nil
	} else {
		resp, err := d.d.Sell(&models.SellParams{
			InstrumentName: ticker,
			Amount:         size,
			Label:          ticker + TRIGGER + "sell",
			Price:          price,
			ReduceOnly:     reduceOnly,
			Type:           "stop_market",
		})
		if err != nil {
			return exchange.TriggerOrder{}, err
		}
		o := resp.Order

		id, _ := strconv.ParseInt(o.OrderID, 10, 64)

		return exchange.TriggerOrder{
			Id:         id,
			Ticker:     o.InstrumentName,
			Side:       false,
			Size:       o.Amount,
			Price:      o.Price.ToFloat64(),
			ReduceOnly: o.ReduceOnly,
			Created:    time.Unix(o.CreationTimestamp/1000, 0),
		}, nil
	}
}
