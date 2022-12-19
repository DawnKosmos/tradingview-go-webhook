package deribit

func (d *Deribit) Cancel(Side int64, Ticker string) error {
	var label string
	switch Side {
	case 1:
		label = Ticker + TRIGGER + "buy"
	case -1:
		label = Ticker + TRIGGER + "sell"
	case 0:
		_, err := d.CancelByLabel(&CancelByLabelRequest{Label: Ticker + "buy"})
		if err != nil {
			return err
		}
		_, err = d.CancelByLabel(&CancelByLabelRequest{Label: Ticker + "sell"})
		return err
	}

	_, err := d.CancelByLabel(&CancelByLabelRequest{Label: label})
	return err
}

func (d *Deribit) CancelTrigger(Side int64, Ticker string) error {
	var label string
	switch Side {
	case 1:
		label = Ticker + "buy"
	case -1:
		label = Ticker + "sell"
	case 0:
		_, err := d.CancelByLabel(&CancelByLabelRequest{Label: Ticker + TRIGGER + "buy"})
		if err != nil {
			return err
		}
		_, err = d.CancelByLabel(&CancelByLabelRequest{Label: Ticker + TRIGGER + "sell"})
		return err
	}

	_, err := d.CancelByLabel(&CancelByLabelRequest{Label: label})
	return err
}

type CancelByLabelRequest struct {
	Label string `json:"label"`
}

func (d *Deribit) CancelByLabel(params *CancelByLabelRequest) (result int, err error) {
	err = d.d.Call("private/cancel_by_label", params, &result)
	return
}
