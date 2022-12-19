package parser

import (
	"fmt"
	"tvalert/cle"
	"tvalert/cle/lexer"
)

type Cancel struct {
	Side         int64
	Ticker       []string
	triggerOrder bool
}

func ParseCancel(tk []lexer.Token) (p Parser, err error) {

	var cancel Cancel
	cancel.Side = 0
	cancel.Ticker = make([]string, 0)

	if len(tk) == 0 {
		return &cancel, nil
	}

	for _, v := range tk {
		switch v.Type {
		case lexer.SIDE:
			if v.Content == "buy" {
				cancel.Side = 1
			} else {
				cancel.Side = -1
			}
		case lexer.FLAG:
			switch v.Content {
			case "stop":
				cancel.triggerOrder = true
			default:
				return nil, nerr(empty, fmt.Sprintf("Error Parsing Cancel, Invalid flag %s", v.Content))
			}
		case lexer.VARIABLE:
			cancel.Ticker = append(cancel.Ticker, v.Content)
		default:
			return nil, nerr(empty, fmt.Sprintf("Error Parsing Cancel, Invalid Type %d %s", v.Type, v.Content))
		}
	}
	return &cancel, nil
}

func (c *Cancel) Evaluate(w cle.Logger, f cle.CLEIO) error {
	if len(c.Ticker) == 0 {
		err := f.Cancel(c.Side, "")
		if err != nil {
			return err
		}
		w.Write([]byte("Orders cancelled succesfully"))
	}
	for _, v := range c.Ticker {
		err := f.Cancel(c.Side, v)
		if err != nil {
			w.ErrorMessage(nerr(err, "Error cancelation"))
			return err
		}
		w.Write([]byte(fmt.Sprintf("%s cancelled succesfully", v)))
	}
	return nil
}
