package listMarketCatalogue

import (
	"testing"
	"time"
	"fmt"
	"sync"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
	"heroku.com/betfairs/aping"
	"strconv"
)

func TestEvent(t *testing.T) {

	const eventID = 28513272
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
			fmt.Println(err)
			return
		}
		event := markets[0].Event

		fmt.Println( "list markets catalogue:", time.Since(startTime) )
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"event", event.Name})
		table.Append([]string{ "date", fmt.Sprintf("%v", event.OpenDate)  })
		table.Append([]string{ "country code", event.CountryCode  })
		table.Append([]string{ "competition", markets[0].Competition.Name  })
		table.Append([]string{ "sport", markets[0].EventType.Name  })
		table.Render()

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"№", "ID", "MARKET NAME", "RUNNER NAME"})
		for i, x := range markets {
			if strings.Contains( x.Name, "Азиатск") || len(x.Runners) == 0 {
				continue
			}
			table.Append([]string{ strconv.Itoa(i), x.ID, x.Name, x.Runners[0].Name  })
			for _,r := range x.Runners[1:] {
				table.Append([]string{"","", "", r.Name, })
			}
		}
		table.Render()

		wg.Done()
	}()
	wg.Wait()
}
var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
)


