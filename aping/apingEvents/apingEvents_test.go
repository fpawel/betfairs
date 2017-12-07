package apingEvents

import (
	"testing"
	"time"
	"fmt"
	"sync"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
	"heroku.com/betfairs/aping"
)

func TestEvent(t *testing.T) {

	const eventID = 28498474 //28495921 //28490335
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

	reader := NewSyncReader(session)
	startTime = time.Now()
	var wg sync.WaitGroup


	for i := 0; i<20; i++ {
		wg.Add(1)
		go func() {
			reader.ReadEvent(eventID)
			wg.Done()
		}()
	}
	go func() {
		wg.Add(1)

		event,err := reader.ReadEvent(eventID)
		if err != nil {
			fmt.Println(err)
		}
		if event == nil {
			return
		}

		fmt.Println( "list markets catalogue:", time.Since(startTime) )
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"event", event.Name})
		table.Append([]string{ "date", fmt.Sprintf("%v", event.OpenDate)  })
		table.Append([]string{ "country code", event.CountryCode  })
		table.Append([]string{ "competition", event.Competition.Name  })
		table.Append([]string{ "sport", event.EventType.Name  })
		table.Render()

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"№", "MARKET NAME", "RUNNER NAME"})
		for i, x := range event.Markets {
			if strings.Contains( x.Name, "Азиатск") || len(x.Runners) == 0 {
				continue
			}
			table.Append([]string{ fmt.Sprintf("%d", i+1), x.Name, x.Runners[0].Name  })
			for _,r := range x.Runners[1:] {
				table.Append([]string{
					"", "", r.Name, })
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

