package main

import (
	"fmt"
	"github.com/fpawel/betfairs/aping/listMarketBook"
	"github.com/fpawel/betfairs/aping/listMarketCatalogue"
	"github.com/fpawel/betfairs/football"
	"sync/atomic"
	"github.com/fpawel/betfairs/football/football2"
	"github.com/fpawel/betfairs/football/football3"
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



