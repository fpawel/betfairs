package football2

import (
	"fmt"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/football"
	"sync/atomic"
)

type Games []Game

var ErrorInterrupted = fmt.Errorf("INTERRUPTED")

type GamesReader struct {
	footballReader      *football.GamesReader
	listMarketCatalogue *listMarketCatalogue.Reader
	listMarketBook      *listMarketBook.Reader
}

func NewGamesReader(f *football.GamesReader, c *listMarketCatalogue.Reader, b *listMarketBook.Reader) (x *GamesReader) {
	return &GamesReader{
		footballReader:      f,
		listMarketCatalogue: c,
		listMarketBook:      b,
	}
}




func (x *GamesReader) Read(interrupt *int32) (games2 Games, err error) {
	var games []football.Game
	games, err = x.footballReader.Read()
	if err != nil {
		return
	}
	if atomic.LoadInt32(interrupt) > 0 {
		err = ErrorInterrupted
		return
	}
	for _, game := range games {
		game := Game{Game: game}
		game.Read(x.listMarketCatalogue, x.listMarketBook)
		if atomic.LoadInt32(interrupt) > 0 {
			err = ErrorInterrupted
			return
		}
		games2 = append(games2, game)
	}
	return
}
