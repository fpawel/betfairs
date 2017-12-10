package aping

type MarketBooks []MarketBook

func (x MarketBooks) Dub() (dub MarketBooks){
	dub = make(MarketBooks, len(x))
	for n := range x {
		dub[n] = x[n].Dub()
	}
	return
}
