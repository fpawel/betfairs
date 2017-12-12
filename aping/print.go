package aping

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"

	"strings"
	"strconv"

)

func PrintMarketCatalogues(markets MarketCatalogues) {
	event := markets[0].Event
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
}


func PrintMarketBook(books MarketBooks, catalogues MarketCatalogues) {
	table := tablewriter.NewWriter(os.Stdout)
	for nMarket, market := range books {
		marketCatalogue,_ := catalogues.Market(market.ID)
		if strings.Contains( marketCatalogue.Name, "Азиатск") {
			continue
		}
		cellMarketName := fmt.Sprintf("%d. %s", nMarket+1, marketCatalogue.Name, )

		for nRunner, runner := range market.Runners {
			runnerCatalogue, _ := catalogues.Runner(market.ID, runner.ID)
			cellRunnerName := fmt.Sprintf("%d. %s", nRunner+1, runnerCatalogue.Name )

			if nRunner > 0 {
				cellMarketName = ""
			}


			b := runner.PriceSize(Back)
			l := runner.PriceSize(Lay)
			n := len(b)
			if len(l) > len(b) {
				n = len(l)
			}
			for i:=0; i<n; i++{
				if i > 0 {
					cellMarketName = ""
					cellRunnerName = ""
				}

				var bP, bS, lP, lS string
				if i< len(b) {
					bP = fmt.Sprintf("%v", b[i].Price)
					bS = fmt.Sprintf("%v $", b[i].Size)
				}
				if i< len(l) {
					lP = fmt.Sprintf("%v", l[i].Price)
					lS = fmt.Sprintf("%v $", l[i].Size)
				}
				table.Append([]string{   cellMarketName, cellRunnerName, strconv.Itoa(i+1), bP, bS, lP, lS  })

			}
		}
	}
	table.Render()
}

