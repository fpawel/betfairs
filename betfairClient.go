package main

import (
	"fmt"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/football"
	"sync/atomic"
	"heroku.com/betfairs/football2"
)

var ErrorInterrupted = fmt.Errorf("INTERRUPTED")

type BetfairClient struct {
	Football            *football.GamesReader
	ListMarketCatalogue *listMarketCatalogue.Reader
	ListMarketBook      *listMarketBook.Reader
}

func (x *BetfairClient) ReadFootballGames2(interrupt *int32) (games2 football2.Games, err error) {
	var games []football.Game
	games, err = x.Football.Read()
	if err != nil {
		return
	}
	if atomic.LoadInt32(interrupt) > 0 {
		err = ErrorInterrupted
		return
	}
	for _, game := range games {
		game := football2.Game{Game: game}
		game.Read(x.ListMarketCatalogue, x.ListMarketBook)
		if atomic.LoadInt32(interrupt) > 0 {
			err = ErrorInterrupted
			return
		}
		games2 = append(games2, game)
	}
	return
}


func (x *BetfairClient) ReadFootballGames3() (games3 []football2.Game3, err error) {
	var games []football.Game
	games, err = x.Football.Read()
	if err != nil {
		return
	}
	for _, game := range games {
		if !game.InPlay {
			continue
		}
		game := football2.Game3{
			Game:game,
		}
		if game.Read(x.ListMarketCatalogue, x.ListMarketBook) == nil {
			games3 = append(games3, game)
		}

	}
	return
}

func (x *BetfairClient) ReadFootballGames4() (games4 []football2.Game4, err error) {
	var games []football.Game
	games, err = x.Football.Read()
	if err != nil {
		return
	}
	for _, game := range games {
		if !game.InPlay {
			continue
		}
		game := football2.Game3{
			Game:game,
		}
		if game.Read(x.ListMarketCatalogue, x.ListMarketBook) == nil {
			games4 = append(games4, game.Game4())
		}

	}
	return
}


