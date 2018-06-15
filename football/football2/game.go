package football2

import (
	"fmt"
	"github.com/fpawel/betfairs/aping"
	"github.com/fpawel/betfairs/aping/listMarketBook"
	"github.com/fpawel/betfairs/aping/listMarketCatalogue"
	"github.com/fpawel/betfairs/countries"
	"github.com/fpawel/betfairs/football"
	"log"
	"math"
	"strings"
	"time"
)

type Games []Game

type Game struct {
	football.Game
	Competition    string    `json:"competition"`
	Country        string    `json:"country"`
	OpenDate       time.Time `json:"open_date"`
	WinBack        float64   `json:"win_back"`
	WinLay         float64   `json:"win_lay"`
	LoseBack       float64   `json:"lose_back"`
	LoseLay        float64   `json:"lose_lay"`
	DrawBack       float64   `json:"draw_back"`
	DrawLay        float64   `json:"draw_lay"`
	TotalMatched   float64   `json:"total_matched"`
	TotalAvailable float64   `json:"total_available"`
	Error          string    `json:"error"`
}

func (x Game) Changes(y Game) (r GameChanges) {
	r.ID = x.ID
	if x.ID != y.ID {
		log.Fatal("IDs must be the same")
	}

	if x.Order != y.Order {
		r.Order = &y.Order
	}
	if x.ScoreHome != y.ScoreHome {
		r.ScoreHome = &y.ScoreHome
	}
	if x.ScoreAway != y.ScoreAway {
		r.ScoreAway = &y.ScoreAway
	}
	if x.InPlay != y.InPlay {
		r.InPlay = &y.InPlay
	}
	if x.Time != y.Time {
		r.Time = &y.Time
	}
	if x.Competition != y.Competition {
		r.Competition = &y.Competition
	}
	if x.Country != y.Country {
		r.Country = &y.Country
	}
	if x.WinBack != y.WinBack {
		r.WinBack = &y.WinBack
	}
	if x.WinLay != y.WinLay {
		r.WinLay = &y.WinLay
	}
	if x.LoseBack != y.LoseBack {
		r.LoseBack = &y.LoseBack
	}
	if x.LoseLay != y.LoseLay {
		r.LoseLay = &y.LoseLay
	}
	if x.DrawBack != y.DrawBack {
		r.DrawBack = &y.DrawBack
	}
	if x.DrawLay != y.DrawLay {
		r.DrawLay = &y.DrawLay
	}
	if x.TotalMatched != y.TotalMatched {
		r.TotalMatched = &y.TotalMatched
	}
	if x.TotalAvailable != y.TotalAvailable {
		r.TotalAvailable = &y.TotalAvailable
	}
	if x.Error != y.Error {
		r.Error = &y.Error
	}
	return
}

func (x *Game) Read(marketCatalogueReader *listMarketCatalogue.Reader, marketBookReader *listMarketBook.Reader) {
	mc, err := marketCatalogueReader.Read(x.ID)
	defer func() {
		if err != nil {
			x.Error = err.Error()
		} else {
			x.Error = ""
		}
	}()

	if err != nil {
		return
	}
	x.OpenDate = mc[0].Event.OpenDate
	x.Competition = mc[0].Competition.Name
	if strings.ToLower(x.Competition) == "чемпионшип" {
		x.Competition = "Чемпионат Футбольной лиги Англии"
	}

	c := countries.ByAlpha2(mc[0].Event.CountryCode)
	if c != nil {
		x.Country = c.Name
	} else {
		x.Country = mc[0].Event.CountryCode
	}

	if strings.Contains(x.Competition, x.Country+" ") {
		x.Competition = strings.Replace(x.Competition, x.Country+" ", " ", -1)
	}

	mainMarket, ok := mc.MainMarket()
	if !ok {
		err = fmt.Errorf("рынок ставок на результат не найден")
		return
	}

	if len(mainMarket.Runners) != 3 {
		err = fmt.Errorf("main market must have 3 runners")
		return
	}

	t := time.Second
	if !x.InPlay {
		t = time.Minute
	}

	var mb aping.MarketBooks
	mb, err = marketBookReader.Read([]aping.MarketID{mainMarket.ID}, t)
	if err != nil {
		return
	}
	for _, m := range mb {
		if m.ID == mainMarket.ID {
			prices6 := m.Prices6()
			x.WinBack, x.WinLay, x.LoseBack, x.LoseLay, x.DrawBack, x.DrawLay = prices6[0], prices6[1], prices6[2], prices6[3], prices6[4], prices6[5]
			x.TotalMatched = math.Round(m.TotalMatched)
			x.TotalAvailable = math.Round(m.TotalAvailable)
			break
		}
	}
}
