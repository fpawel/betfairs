package main

import (
	"fmt"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/football"
	"sync/atomic"
	"heroku.com/betfairs/football/football2"
	"heroku.com/betfairs/football/football3"
	"heroku.com/betfairs/football/football4"
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


func (x *BetfairClient) ReadFootballGames3(interrupt *int32) (games3 []football3.Game, err error) {
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
		if !game.InPlay {
			continue
		}
		game := football3.Game{
			Game:game,
		}
		if game.Read(x.ListMarketCatalogue, x.ListMarketBook) == nil {
			games3 = append(games3, game)
		}
		if atomic.LoadInt32(interrupt) > 0 {
			err = ErrorInterrupted
			return
		}
	}
	return
}

func (x *BetfairClient) ReadFootballGames4(interrupt *int32) (games4 []football4.Game, err error) {
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
		if !game.InPlay {
			continue
		}
		var game4 football4.Game
		game4,err = football4.ReadGame(game, x.ListMarketCatalogue, x.ListMarketBook)
		if err == nil {
			games4 = append(games4, game4)
		}
		if atomic.LoadInt32(interrupt) > 0 {
			err = ErrorInterrupted
			return
		}
	}
	return
}


