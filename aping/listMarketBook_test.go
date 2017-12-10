package aping

import (
	"testing"
	"time"
	"fmt"
	"os"

)

func TestEventPrices(t *testing.T) {

	session :=  NewSession(adminBetfairUser, adminBetfairPass)
	{
		session := session.GetSession()
		fmt.Println(session.SessionToken, session.AppKey)
	}

	marketCatalogues,err := session.ListMarketCatalogue(28494326)
	if err != nil {
		t.Error(err)
		return
	}

	var marketBooks MarketBooks
	for _,marketIDs := range marketCatalogues.Take40MarketIDs() {
		marketBook,err := session.ListMarketBook(marketIDs)
		if err != nil {
			t.Error(err)
			return
		}
		marketBooks = append(marketBooks, marketBook ...)
	}
	PrintMarketBook(marketBooks,marketCatalogues)
}

func formatTime(t time.Time) string {
	if t.Before( time.Unix(1,1) ) {
		return ""
	}
	return t.Add(time.Hour * 3).Format("02.01.06 15:04")

}

var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
)
