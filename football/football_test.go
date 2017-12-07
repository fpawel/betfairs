package football

import (
	"testing"
	"time"
	"fmt"

)

func TestFetchGames(t *testing.T) {
	startTime := time.Now()
	games, err := FetchGames()
	fmt.Println("fetch", time.Since(startTime), err)
	PrintGames(games)
}

func TestCache(t *testing.T) {
	cache := new(SyncReader)

	startTime := time.Now()
	for i:=0; i<20; i++ {
		i := i
		go func() {
			games,_ := cache.Read()
			fmt.Println(i, ":", len(games))
		}()
	}
	games,err := cache.Read()
	fmt.Println("fetch", time.Since(startTime), err)
	startTime = time.Now()
	PrintGames(games)
}


