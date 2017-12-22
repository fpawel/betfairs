package aping


import (
	"sync"
	"time"
	"encoding/json"
	"fmt"

)

type Session struct {
	user, pass	string
	muConsumers, muSessionToken sync.RWMutex
	consumers                   []chan SessionResult
	sessionToken, appKey        string
	sessionTime                 time.Time
}

type SessionResult struct {
	SessionToken, AppKey string

	Error error
}

func NewSession(user, pass	string) *Session {
	return &Session{
		user: user,
		pass:pass,
	}
}

func (x *Session) GetSession() SessionResult {

	x.muSessionToken.RLock()
	if x.appKey != "" && x.sessionToken != "" && time.Since(x.sessionTime) < 30*time.Minute {
		x.muSessionToken.RUnlock()
		return SessionResult{SessionToken: x.sessionToken, AppKey:x.appKey, Error: nil}
	}
	x.muSessionToken.RUnlock()


	resultChan := make(chan SessionResult)

	x.muConsumers.Lock()
	x.consumers = append(x.consumers, resultChan)


	if len(x.consumers) == 1 {
		go func() {
			result := SessionResult{}
			result.SessionToken,result.Error = Login(x.user,x.pass)

			if result.Error == nil {
				x.muSessionToken.Lock()
				x.sessionToken = result.SessionToken
				x.sessionTime = time.Now()
				x.muSessionToken.Unlock()
				result.AppKey, result.Error = fetchAppKey(x.sessionToken)
				if result.Error == nil {
					x.appKey = result.AppKey
				}
			}

			x.muConsumers.Lock()
			defer x.muConsumers.Unlock()
			for _, ch := range x.consumers {
				ch <- result
			}
			x.consumers = nil
		}()
	}

	x.muConsumers.Unlock()

	return <-resultChan
}


func (x *Session) getResponse(endpoint Endpoint, params interface{}) (responseBody []byte, err error) {
	r := x.GetSession()
	if r.Error != nil {
		return nil, r.Error
	}
	return endpoint.getResponse(r.SessionToken, &r.AppKey, params)
}

//ReadEvent - получить каталог события
func (x *Session) ListMarketCatalogue(eventID int) (MarketCatalogues, error) {

	type listMarketCatalogueRequest struct {
		Locale           string       `json:"locale"`
		MarketProjection []string     `json:"marketProjection"`
		MarketFilter     MarketFilter `json:"filter"`
		MaxResults       int          `json:"maxResults"`
	}

	request := &listMarketCatalogueRequest{
		Locale: "ru",
		MarketFilter:MarketFilter{
			EventIDs: []int{eventID},
		},
		MarketProjection: []string{"EVENT", "EVENT_TYPE", "COMPETITION", "RUNNER_DESCRIPTION"},
		MaxResults: 1000,
	}


	responseBody, err := x.getResponse(BettingAPIEndpoint("listMarketCatalogue"), &request)
	if err != nil {
		return nil, err
	}

	var response struct {
		Result MarketCatalogues `json:"result"`
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("%q, %q", err, string(responseBody))
	}

	if len(response.Result) == 0{
		return nil, fmt.Errorf("no markets in catalogue: %s", string(responseBody))
	}

	return response.Result, nil
}



func (x *Session) ListMarketBook(marketIDs []MarketID) (MarketBooks, error) {

	type readMarketBookRequest struct {
		Locale          string   `json:"locale"`
		MarketIDs         []MarketID `json:"marketIds"`
		PriceProjection PriceProjection `json:"priceProjection"`
	}

	request := readMarketBookRequest{
		Locale : "ru",
		MarketIDs : marketIDs,
		PriceProjection: PriceProjection{
			PriceData: []string{"EX_BEST_OFFERS"},
			Virtualise: true,
		},
	}
	responseBody, err := x.getResponse(BettingAPIEndpoint("listMarketBook"), &request)
	if err != nil {
		return nil,err
	}
	var r struct {
		Result MarketBooks `json:"result"`
	}
	err = json.Unmarshal(responseBody, &r)
	if err != nil {
		return nil, fmt.Errorf("%q, %q", err, string(responseBody))
	}
	if len(r.Result) == 0{
		return nil, fmt.Errorf("no markets in book: %s", string(responseBody))
	}
	return r.Result,err
}


