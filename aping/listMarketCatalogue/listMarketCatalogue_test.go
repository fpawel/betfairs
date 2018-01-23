package listMarketCatalogue

import (
	"testing"
	"time"
	"fmt"
	"sync"
	"os"
	"github.com/fpawel/betfairs/aping"
	"strconv"
)

func TestEvent(t *testing.T) {

	const eventID = 28515945
	// 28490325

	//startTime := time.Now()
	//competition, err := fetchFootballCompetitionInfo(eventID, true, )
	//fmt.Println("competition:", time.Since(startTime), err)

	startTime := time.Now()
	session :=  aping.NewSession(adminBetfairUser, adminBetfairPass)
	{
		session := session.GetSession()
		fmt.Println(session.SessionToken, session.AppKey, time.Since(startTime))
	}

	reader := New(session)
	startTime = time.Now()
	var wg sync.WaitGroup


	for i := 0; i<20; i++ {
		i := i
		wg.Add(1)
		go func() {
			tmp,err := reader.Read(eventID)
			if err == nil {
				tmp[0].Name = strconv.Itoa(i)
				tmp[0].Runners[0].Name = strconv.Itoa(i)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Add(1)

		markets,err := reader.Read(eventID)
		if err != nil {
			t.Error(err)
			wg.Done()
			return
		}
		aping.PrintMarketCatalogues(markets)
		wg.Done()
	}()
	wg.Wait()
}
var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
)


