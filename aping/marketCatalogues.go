package aping


type MarketCatalogues []MarketCatalogue

func (x MarketCatalogues) Dub() (dub MarketCatalogues){
	dub = make(MarketCatalogues, len(x))
	for n := range x {
		dub[n] = x[n].Dub()
	}
	return
}


func (x MarketCatalogues) Take40MarketIDs() (r [][] string){
	for i,y := range x {
		if i % 40 == 0 {
			r = append(r, []string{})
		}
		n := len(r)-1
		r[n] = append(r[n], y.ID)
	}
	return
}

func (x MarketCatalogues) MainMarket() (m MarketCatalogue, ok bool){
	for n := range x {
		if x[n].Name == "Ставки на результат" {
			ok = true
			m = x[n]
			break
		}
	}
	return
}

func (x MarketCatalogues) Market(id string) (m MarketCatalogue, ok bool){
	for n := range x {
		if x[n].ID == id {
			ok = true
			m = x[n]
			break
		}
	}
	return
}

func (x MarketCatalogues) Runner(marketID string, runnerID int) (r RunnerCatalogue, ok bool){
	var m MarketCatalogue
	m,ok = x.Market(marketID)
	if ok {
		r,ok = m.Runner(runnerID)
	}
	return
}