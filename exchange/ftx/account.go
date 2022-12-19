package ftx

import "tradingviewListener/exchange"

func (p *FTX) OpenPosition() (out []exchange.Position, err error) {
	var ff PositionsResponse
	resp, err := p.get("positions?showAvgPrice=true", []byte(""))
	if err != nil {
		return
	}
	err = processResponse(resp, &ff)
	if err != nil {
		return
	}
	for _, v := range ff.Result {
		if v.NotionalSize > 0.05 || v.NotionalSize < -0.05 {
			out = append(out, v)
		}
	}
	return out, nil
}

func (p *FTX) Account() (exchange.Account, error) {
	var ff AccountResponse
	var out exchange.Account
	resp, err := p.get("account", []byte(""))
	if err != nil {
		return out, err
	}
	err = processResponse(resp, &ff)
	if err != nil {
		return out, err
	}
	out = ff.Result
	var pp []exchange.Position

	for _, v := range out.Positions {
		if v.NotionalSize > 0.1 || v.NotionalSize < -0.1 {
			pp = append(pp, v)
		}
	}
	out.Positions = pp
	return out, nil
}
