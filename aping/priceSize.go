package aping

type PriceSize struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type PriceSizes []PriceSize

func (x PriceSizes) Dub() (dub PriceSizes) {
	dub = make([]PriceSize, len(x))
	for i:=range x{
		dub[i] = x[i]
	}
	return
}

