package aping


type ExchangePrices struct {
	AvailableToBack PriceSizes `json:"availableToBack,omitempty"`
	AvailableToLay  PriceSizes `json:"availableToLay,omitempty"`
	TradedVolume    PriceSizes `json:"tradedVolume,omitempty"`
}

func (x ExchangePrices) Back() (b float64){
	for _,v := range x.AvailableToBack {
		if v.Price > b {
			b = v.Price
		}
	}
	return
}

func (x ExchangePrices) Lay() (l float64){
	l = 1000
	for _,v := range x.AvailableToLay {
		if v.Price < l {
			l = v.Price
		}
	}
	return
}

func (x ExchangePrices) Dub() (dub ExchangePrices)  {
	dub.AvailableToBack = x.AvailableToBack.Dub()
	dub.AvailableToLay = x.AvailableToLay.Dub()
	dub.TradedVolume = x.TradedVolume.Dub()
	return
}
