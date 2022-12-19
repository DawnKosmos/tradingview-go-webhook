package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

const (
	VARIABLE     TokenType = iota //x
	TICKER                        // btc-perp
	SIDE                          // buy, sell
	STOP                          // stop
	FLOAT                         // 100 => 100 $ of btc
	UFLOAT                        // u100  u = unitFloat => buying 100 btc
	PERCENT                       //
	DFLOAT                        // -200 d = differenceFloat =>
	ASSIGN                        // =
	FLAG                          // -l -le
	DELETE                        // delete a b c
	FUNC                          // func(a,b,c) creating function
	DURATION                      // 4h 1d 30
	LBRACKET                      // (
	RBRACKET                      // )
	MARKET                        // -market
	SOURCE                        // -high -low -open -close
	CANCEL                        // cancel
	CLOSE                         // close
	FUNDINGPAYS                   // fpay | fundingpayments
	POSITION                      // -position
	FUNDINGRATES                  //fundingrates
	ACCOUNT
	PRICEORDER
)

type Token struct {
	Type    TokenType
	Content string
}

type lexerError struct {
	input      string
	err        error
	errmessage string
}

func nerr(input string, err error, errmessage string) *lexerError {
	return &lexerError{input, err, errmessage}
}

func (l *lexerError) Error() string {
	return fmt.Sprintf("Text: %s, Error: %v + %s", l.input, l.err, l.errmessage)
}

// Lexer converts an input to a Token Array
func Lexer(input string) (t []Token, err error) {
	in := strings.Split(input, " ") //tokens are seperated with whitespaces. Only exeptions are function, here its forbidden to seperate

	for _, s := range in {
		if len(s) == 0 {
			continue
		}
		last := len(s) - 1
		switch s {
		case "buy", "sell":
			t = append(t, Token{SIDE, s})
		case "stop":
			t = append(t, Token{STOP, s})
		case "delete":
			t = append(t, Token{DELETE, "delete"})
		case "=":
			t = append(t, Token{ASSIGN, "="})
		case "cancel":
			t = append(t, Token{CANCEL, "cancel"})
		case "fpays", "fundingpays":
			t = append(t, Token{FUNDINGPAYS, "fpays"})
		case "frates", "fundingrates":
			t = append(t, Token{FUNDINGRATES, "frates"})
		case "account", "acc":
			t = append(t, Token{ACCOUNT, "account"})
		case "close":
			t = append(t, Token{CLOSE, s})
		default:
			if (s[last] == 'h' || s[last] == 'm' || s[last] == 'd') && len(s) > 1 {
				_, err := strconv.Atoi(s[:last])
				if err == nil {
					t = append(t, Token{DURATION, s})
					continue
				}
			}

			if len(s) > 6 {
				if s[:5] == "func(" {
					t = append(t, Token{FUNC, "func"}, Token{LBRACKET, "("})
					t = append(t, lexFunc([]byte(s[5:]))...)
					continue
				}
			}

			if s[0] == '-' {
				_, err := strconv.ParseFloat(s[1:], 64)

				if err == nil {
					t = append(t, Token{DFLOAT, s[1:]})
				} else {
					ss := s[1:]
					switch ss {
					case "po":
						t = append(t, Token{PRICEORDER, "1.0"})
					case "low", "high", "open", "close":
						t = append(t, Token{SOURCE, ss})
					case "position":
						t = append(t, Token{POSITION, "1.0"})
					case "market":
						t = append(t, Token{MARKET, "1.0"})
					default:
						t = append(t, Token{FLAG, ss})
					}
				}
				continue
			}

			if s[0] == 'u' && len(s) > 1 {
				_, err := strconv.Atoi(s[1:])
				if err == nil {
					t = append(t, Token{UFLOAT, s[1:]})
				} else {
					t = append(t, Token{VARIABLE, s})
				}
				continue
			}

			if s[last] == '%' {
				_, err := strconv.ParseFloat(s[:last], 64)
				if err != nil {
					return t, nerr(s, err, "A variable can't end with %")
				}
				t = append(t, Token{PERCENT, s[:len(s)-1]})
				continue
			}

			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				t = append(t, Token{FLOAT, s})
				continue
			}
			t = append(t, lexVariable([]byte(s))...)
		}
	}

	return
}

func lexFunc(s []byte) []Token {
	var temp []byte
	var tk []Token
	for _, v := range s {
		switch v {
		case ')':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			tk = append(tk, Token{VARIABLE, string(temp)})
			temp = []byte("")
		default:
			temp = append(temp, v)
		}
	}

	return tk
}

// LexVariable lexes functions e.g. a(xrp-buy,5,10)
func lexVariable(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{LBRACKET, ""})
			temp = []byte("")
		case ')':
			l, _ := Lexer(string(temp))
			tk = append(tk, l...)
			tk = append(tk, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			temp = append(temp, ' ')
		default:
			temp = append(temp, v)
		}
	}

	if len(temp) != 0 {
		tk = append(tk, Token{VARIABLE, string(temp)})
	}
	return tk
}

func (t Token) Stringer() (out string) {

	switch t.Type {
	case VARIABLE:
		return t.Content
	case CANCEL:
		out = "cancel"
	case UFLOAT:
		out = "u" + t.Content
	case PERCENT:
		out = t.Content + "%"
	case FLAG:
		out = "-" + t.Content
	case ASSIGN:
		out = "="
	case TICKER:
		return t.Content
	case SIDE:
		return t.Content
	case STOP:
		return "stop"
	case FLOAT:
		return t.Content
	case DFLOAT:
		return "-" + t.Content
	case FUNC:
		return t.Content
	case LBRACKET:
		return "("
	case RBRACKET:
		return ")"
	case SOURCE:
		return "-" + t.Content
	case CLOSE:
		return "close"
	case FUNDINGPAYS:
		return "fpays"
	case POSITION:
		return "-position"
	case FUNDINGRATES:
		return "frates"
	case PRICEORDER:
		return "-po"
	default:
		return t.Content
	}

	return
}
