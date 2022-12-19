package deribit

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
)

func TestDeribit(t *testing.T) {
	n, err := os.ReadFile("private.key")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	var name, public, private string
	_, err = fmt.Fscanf(bytes.NewReader(n), "%s %s %s", &name, &public, &private)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	d, err := NewPrivate(context.Background(), name, public, private)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	res, err := d.FreeCollateral("BTC-Perpetual")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(res)
}
