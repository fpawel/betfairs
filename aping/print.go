package aping

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"

	"strings"
	"strconv"
)

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

