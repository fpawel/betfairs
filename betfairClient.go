package main

import (
	"fmt"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/football"
	"sync/atomic"
	"heroku.com/betfairs/football2"
	"time"
	"heroku.com/betfairs/aping"
)

var ErrorInterrupted = fmt.Errorf("INTERRUPTED")

type BetfairClient struct {
	Football            *football.GamesReader
	ListMarketCatalogue *listMarketCatalogue.Reader
	ListMarketBook      *listMarketBook.Reader
}

func (x *BetfairClient) ReadFootballGames(interrupt *int32) (games2 football2.Games, err error) {
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


func (x *BetfairClient) ReadFootballPrices() error {
	games, err := x.Football.Read()
	if err != nil {
		return err
	}
	for _,game := range games {
		if !game.InPlay {
			continue
		}
		marketCatalogues, err := x.ListMarketCatalogue.Read(game.ID)
		if err != nil {
			return err
		}
		var marketBooks aping.MarketBooks
		for _,xs := range marketCatalogues.Take40MarketIDs(){
			ms,err := x.ListMarketBook.Read(xs, time.Hour)
			if err != nil {
				return err
			}
			marketBooks = append(marketBooks, ms ...)
		}
	}







}