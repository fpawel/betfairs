package aping


type ExchangePrices struct {
	AvailableToBack PriceSizes `json:"availableToBack"`
	AvailableToLay  PriceSizes `json:"availableToLay"`
	TradedVolume    PriceSizes `json:"tradedVolume"`
}

func (x ExchangePrices) Back() (b float64){
	if len(x.AvailableToBack)>0 {
		b = x.AvailableToBack[0].Price
	}
	return
}

func (x ExchangePrices) Lay() (l float64){
	if len(x.AvailableToLay)>0 {
		l = x.AvailableToLay[0].Price
	}
	return
}

func (x ExchangePrices) Dub() (dub ExchangePrices)  {
	dub.AvailableToBack = x.AvailableToBack.Dub()
	dub.AvailableToLay = x.AvailableToLay.Dub()
	dub.TradedVolume = x.TradedVolume.Dub()
	return
}
