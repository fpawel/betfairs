package football

import (
	"fmt"
	"os"
	"strconv"
	tablewriter "github.com/olekukonko/tablewriter"
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

func (x Game) String() string {
	strScore := ""
	if x.InPlay {
		strScore = fmt.Sprintf(" %d - %d", x.ScoreHome, x.ScoreAway)
	}
	return fmt.Sprintf("%d %s - %s %s%s", x.ID, x.Home, x.Away, x.Time, strScore)
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



