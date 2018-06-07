package football

import (
	"fmt"
	"os"
	"strconv"
	tablewriter "github.com/olekukonko/tablewriter"
	"regexp"
	"time"
)

type Game struct {
	ID        int    `json:"id"`
	Order 	  int  `json:"order"`
	Home      string `json:"home"`
	Away      string `json:"away"`
	ScoreHome int    `json:"score_home"`
	ScoreAway int    `json:"score_away"`
	InPlay    bool   `json:"in_play"`
	Time      string `json:"time"`
}

type GameLive struct {
	ID        int `json:"id"`
	OpenDate time.Time `json:"open_date"`
	ScoreHome int `json:"score_home"`
	ScoreAway int `json:"score_away"`
	Minute    int `json:"minute"`
}


func (x Game) GameLive(openDate time.Time ) GameLive {
	minute,_ := x.Minute()
	return GameLive{
		ID:x.ID,
		OpenDate:openDate,
		ScoreHome:x.ScoreHome,
		ScoreAway:x.ScoreAway,
		Minute:minute,
	}
}

func (x Game) String() string {
	strScore := ""
	if x.InPlay {
		strScore = fmt.Sprintf(" %d - %d", x.ScoreHome, x.ScoreAway)
	}
	return fmt.Sprintf("%d %s - %s %s%s", x.ID, x.Home, x.Away, x.Time, strScore)
}

func (x Game) Minute() (int,error) {
	m := regexp.MustCompile("(\\d+)′").FindStringSubmatch(x.Time)
	if len(m) == 2 {
		return strconv.Atoi(m[1])
	}
	return 0,fmt.Errorf("%s: time does not match", x.Time)
}

func (x Game) HasMinute() bool {
	_,err := x.Minute()
	return err == nil
}

func PrintGames(games []Game){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"№", "ID", "HOME", "AWAY", "SCORE", })
	for i, x := range games {
		strScore := ""
		if x.InPlay {
			strScore = fmt.Sprintf("%d-%d", x.ScoreHome, x.ScoreAway)
		}
		table.Append([]string{strconv.Itoa(i + 1), strconv.Itoa(x.ID), x.Home, x.Away, strScore, x.Time})
	}
	table.Render()
}