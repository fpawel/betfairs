package football

import (
	"strings"
	"os"
	"strconv"
	"io/ioutil"
	"fmt"
	"heroku.com/betfairs/webclient"
	"github.com/PuerkitoBio/goquery"

	"errors"
)

var ErrorNoGames = errors.New("NO GAMES")
var ErrorNotReady = errors.New("NOT READY")
var ErrorHTMLSelectionNotFound = errors.New("NOT FOUND HTML SELECTION")
var ErrorHTMLSelectionDublikat = errors.New("NOT FOUND HTML SELECTION")


func parseGame ( node *goquery.Selection) (Game, error) {

	var (
		x Game
		err error)

	strDataEventID, _ := node.Parent().Parent().Parent().Attr("data-eventid")
	x.ID, err = strconv.Atoi(strDataEventID)
	if err != nil {
		return x, fmt.Errorf("data-eventid not ok: %v", err)

	}

	x.Home = strings.TrimSpace(node.Find("span.home-team-name").Text())
	x.Away = strings.TrimSpace(node.Find("span.away-team-name").Text())

	x.ScoreHome, err = strconv.Atoi(strings.TrimSpace(node.Find("span.ui-score-home").Text()))
	if err == nil {
		x.ScoreAway, err = strconv.Atoi(strings.TrimSpace(node.Find("span.ui-score-away").Text()))
		if err == nil {
			x.InPlay = true
		}
	}

	x.InPlayTime = strings.TrimSpace(node.Find("span.inplay").Text())
	if x.InPlayTime == "" {
		x.InPlayTime = strings.TrimSpace(node.Find("span.date").Text())
	}

	return x, nil
}

func parseGames (document *goquery.Document) ( games []Game, err error) {

	document.Find("div[data-eventid] div.details-event div a").Each(func(i int, node *goquery.Selection) {
		var x Game
		x, err = parseGame(node)
		if err == nil {
			games = append(games, x)
		}
	})

	if len(games) == 0{
		htmlStr,_ := document.Html()
		ioutil.WriteFile("nogames.html", []byte(htmlStr), os.ModeExclusive)
		err = ErrorNoGames
	}
	return
}



func FetchGames() (games []Game, err error) {

	var URLStr string

	if strings.ToLower(os.Getenv("BETFAIR_COM_NOT_ALLOWED") ) == "true" {
		URLStr = "https://betfairs.herokuapp.com/redirect-betfair/sport/football"
	} else {
		URLStr = webclient.NewURL("sport/football")
	}

	err = webclient.Fetch(URLStr, func(document *goquery.Document) error {
		games, err = parseGames(document)
		return err
	})
	return
}