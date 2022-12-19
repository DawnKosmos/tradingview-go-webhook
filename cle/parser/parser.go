package parser

import (
	"errors"
	"fmt"
	"tvalert/cle"
	"tvalert/cle/lexer"
)

type Parser interface {
	Evaluate(w cle.Logger, f cle.CLEIO) error
}

// Parse returns a Parser which then gets Evaluated and returns
func Parse(tk []lexer.Token, c cle.Logger) (Parser, error) {
	nk := tk

	if len(tk) == 0 {
		return nil, nerr(empty, "Error nothing got lexed")
	}

	var o Parser
	var err error

	switch nk[0].Type {
	case lexer.SIDE: // buy, sell
		o, err = ParseOrder(nk[0].Content, nk[1:])
	case lexer.STOP: //stop
		o, err = ParseStop(nk[1:])
	case lexer.CANCEL: //cancel
		o, err = ParseCancel(nk[1:])
	case lexer.CLOSE: //fclose
		o, err = ParseClose(nk[1:])
	default:
		return o, nerr(empty, fmt.Sprintf("Invalid Type Error during Parsing %v", nk[0].Type))
	}

	if err != nil {
		return o, err
	}

	return o, nil
}

type parseError struct {
	err error
	msg string
}

var empty = errors.New("")

func nerr(err error, msg string) *parseError {
	return &parseError{err, msg}
}

func (e *parseError) Error() string {
	return fmt.Sprintf("Message:%s : %v", e.msg, e.err)
}
