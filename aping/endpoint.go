package aping


import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
)

type Endpoint struct {
	URL    string
	Method string
}


func AccauntAPIEndpoint(s string) (r Endpoint) {
	r.URL = "https://api.betfair.com/exchange/account/json-rpc/v1"
	r.Method = "AccountAPING/v1.0/" + s
	return
}

func BettingAPIEndpoint(s string) (r Endpoint) {
	r.URL = "https://api.betfair.com/exchange/betting/json-rpc/v1"
	r.Method = "SportsAPING/v1.0/" + s
	return

}

func (x Endpoint) getResponse( sessionToken string, appKey  *string, params interface{}) (responseBody []byte, err error) {

	jsonReq := struct {
		Jsonrpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
		Id      int         `json:"id"`
	}{"2.0", x.Method, params, 1}

	var reqbytes []byte
	if reqbytes, err = json.Marshal(&jsonReq); err != nil {
		return
	}

	var req *http.Request
	if req, err = http.NewRequest("POST", x.URL, bytes.NewBuffer(reqbytes)); err != nil {
		return
	}
	req.ContentLength = int64(len(reqbytes))
	if appKey != nil {
		req.Header.Set("X-Application", *appKey)
	}

	req.Header.Set("X-Authentication", sessionToken)
	req.Header.Set("ContentType", "application/json")
	req.Header.Set("AcceptCharset", "UTF-8")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err = ioutil.ReadAll(resp.Body)



	return
}
