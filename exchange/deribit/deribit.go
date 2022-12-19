package deribit

import (
	"context"
	"github.com/frankrap/deribit-api"
)

type Deribit struct {
	name string
	d    *deribit.Client
}

func (d *Deribit) Name() string {
	return "Deribit" + d.name
}

func NewPrivate(ctx context.Context, name, apiId, apiSecret string) (*Deribit, error) {
	d := deribit.New(&deribit.Configuration{
		Ctx:           ctx,
		Addr:          deribit.RealBaseURL,
		ApiKey:        apiId,
		SecretKey:     apiSecret,
		AutoReconnect: true,
		DebugMode:     true,
	})

	_, err := d.Test()
	if err != nil {
		return nil, err
	}
	return &Deribit{name: name, d: d}, nil
}
