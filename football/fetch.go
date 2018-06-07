package football

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fpawel/betfairs/webclient"
	"os"
	"strconv"
	"strings"

	"errors"
	"io/ioutil"
	"log"
	"regexp"
)

var ErrorNoGames = errors.New("NO GAMES")
var ErrorNotReady = errors.New("NOT READY")
var ErrorHTMLSelectionNotFound = errors.New("NOT FOUND HTML SELECTION")
var ErrorHTMLSelectionDublikat = errors.New("NOT HTML SELECTION DUBLIKAT")

func parseGame(node *goquery.Selection) (Game, error) {

	var (
		x   Game
		err error
	)

	strDataEventID, _ := node.Attr("data-eventid")
	x.ID, err = strconv.Atoi(strDataEventID)
	if err != nil {
		return x, fmt.Errorf("data-eventid not ok: %v", err)

	}

	x.Home = strings.TrimSpace(node.Find("div.teams-container span.team-name:nth-child(1)").Text())
	x.Away = strings.TrimSpace(node.Find("div.teams-container span.team-name:nth-child(2)").Text())

	if len(strings.TrimSpace(x.Home)) == 0 || len(strings.TrimSpace(x.Away)) == 0{
		x.Home = strings.TrimSpace(node.Find("span.event-runner1").Text())
		x.Away = strings.TrimSpace(node.Find("span.event-runner2").Text())
		x.Time = regexp.MustCompile(`\d\d:\d\d`).FindString(node.Text())
	} else {
		x.ScoreHome, err = strconv.Atoi(strings.TrimSpace(node.Find("span.ui-score-home").Text()))
		if err == nil {
			x.ScoreAway, err = strconv.Atoi(strings.TrimSpace(node.Find("span.ui-score-away").Text()))
			if err == nil {
				x.InPlay = true
			}
		}
		x.Time = strings.TrimSpace(node.Find("span.inplay").Text())
		if x.Time == "" {
			x.Time = strings.TrimSpace(node.Find("span.date").Text())
		}
		x.Time = strings.Replace(x.Time, " (In-Play)", "", 1)
	}

	if len(strings.TrimSpace(x.Home)) == 0{
		return x, errors.New("home team not found")
	}
	if len(strings.TrimSpace(x.Away)) == 0{
		return x, errors.New("away team not found")
	}

	return x, nil
}

func parseGames(document *goquery.Document) (games []Game, err error) {

	document.Find("div[data-eventid]").Each(func(i int, node *goquery.Selection) {
		var x Game
		x, err = parseGame(node)
		x.Order = i
		if err == nil {
			games = append(games, x)
		} else {
			s,_ := node.Html()
			log.Println(err, ":", strings.TrimSpace(s))
		}
	})

	if len(games) == 0 {
		err = ErrorNoGames
	}
	return
}

func FetchGames() (games []Game, err error) {

	var URLStr string

	if len(os.Getenv("LOCALHOST")) > 0 {
		URLStr = "https://betfairs.herokuapp.com/redirect-betfair/sport/football"
	} else {
		URLStr = webclient.NewURL("sport/football")
	}

	err = webclient.Fetch(URLStr, func(document *goquery.Document) error {
		games, err = parseGames(document)
		if err != nil {
			s, err := document.Html()
			if err != nil {
				s = err.Error()
			}
			ioutil.WriteFile("assets/error.html", []byte(s), 0644)
		}
		return err
	})
	return
}
