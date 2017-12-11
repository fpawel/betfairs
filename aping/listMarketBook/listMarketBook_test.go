package listMarketBook

import (
	"os"
	"testing"
	"time"
	"sync"
	"heroku.com/betfairs/aping"
)



func TestEventPrices(t *testing.T) {

	const eventID = 28500450

	session :=  aping.NewSession(adminBetfairUser, adminBetfairPass)
	if err := session.GetSession().Error; err != nil {
		t.Fatal(err)
		return
	}

	marketCalogues,err := session.ListMarketCatalogue(eventID)
	if err != nil {
		t.Fatal(err)
		return
	}
	reader := New(session)
	var wg sync.WaitGroup
	for i := 0; i<20; i++ {
		wg.Add(1)
		go func() {
			for _,xs := range marketCalogues.Take40MarketIDs(){
				tmp,err := reader.Read(xs, time.Hour)
				if err == nil {
					tmp[0].ID = "sdfsdf"
					tmp[0].Runners[0].ID = 0
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Add(1)
		var marketBooks aping.MarketBooks
		for _,xs := range marketCalogues.Take40MarketIDs(){
			ms,err := reader.Read(xs, time.Hour)
			if err != nil {
				t.Fatal(err)
				wg.Done()
				return
			}
			marketBooks = append(marketBooks, ms ...)
		}
		aping.PrintMarketBook(marketBooks,marketCalogues)

		wg.Done()
	}()
	wg.Wait()
}


var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
)

