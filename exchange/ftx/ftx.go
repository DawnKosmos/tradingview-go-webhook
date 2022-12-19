package ftx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const FTXURL = "https://ftx.com/api/"

type FTX struct {
	Client     *http.Client
	key        string
	secret     []byte
	Subaccount string
}

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

func New(cl *http.Client, key, secret, sub string) *FTX {
	c := &FTX{cl, key, []byte(secret), sub}
	return c
}

func (p *FTX) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := p.sign(signaturePayload)
	req, _ := http.NewRequest(method, FTXURL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", p.key)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if p.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", p.Subaccount)
	}
	return req
}

func SignWebhook(key []byte, tnow int64) string {
	signature := fmt.Sprintf("%dwebsocket_login", tnow)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(signature))
	return hex.EncodeToString(mac.Sum(nil))
}

func (p *FTX) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, p.secret)
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (p *FTX) get(path string, body []byte) (*http.Response, error) {
	preparedRequest := p.signRequest("GET", path, body)
	resp, err := p.Client.Do(preparedRequest)
	return resp, err
}

func (p *FTX) post(path string, body []byte) (*http.Response, error) {
	preparedRequest := p.signRequest("POST", path, body)
	resp, err := p.Client.Do(preparedRequest)
	return resp, err
}

func (p *FTX) delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := p.signRequest("DELETE", path, body)
	resp, err := p.Client.Do(preparedRequest)
	return resp, err
}

func processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	return nil
}
