package football2

import (
	"fmt"
	"sync/atomic"
	"heroku.com/betfairs/football"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/aping/listMarketBook"
)

type Games []Game

var ErrorInterrupted = fmt.Errorf("INTERRUPTED")


type GamesReader struct {
	FootballReader *football.GamesReader
	ListMarketCatalogue *listMarketCatalogue.Reader
	ListMarketBook *listMarketBook.Reader
	DoneFlag int32
}

func (x *GamesReader)Read() (games2 Games, err error) {
	var games []football.Game
	games, err = x.FootballReader.Read()
	if err != nil {
		return
	}
	if atomic.LoadInt32(&x.DoneFlag) > 0 {
		return
	}
	for _,game := range games {
		game := Game{Game:game}
		game.Read(x.ListMarketCatalogue,x.ListMarketBook)
		if atomic.LoadInt32(&x.DoneFlag) > 0 {
			return
		}
		games2 = append(games2, game)
	}
	return
}
