package football2

import (
	"time"

	"heroku.com/betfairs/aping"
)

type Game4 struct {
	ID              int       `json:"id"`
	OpenDate        time.Time `json:"openDate"`
	CountryCode     string    `json:"countryCode"`
	Home            string    `json:"home"`
	Away            string    `json:"away"`
	ScoreHome       int       `json:"score_home"`
	ScoreAway       int       `json:"score_away"`
	Minute          *int      `json:"minute"`
	CompetitionID   int       `json:"competition_id"`
	CompetitionName string    `json:"competition_name"`

	Markets []Market4 `json:"markets"`
}

type Market4 struct {
	ID int `json:"marketId"`
	Name string `json:"marketName,omitempty"`
	TotalMatched float64 `json:"totalMatched,omitempty"`
	TotalAvailable float64 `json:"totalAvailable,omitempty"`
	Runners []Runner4 `json:"runners,omitempty"`
}

type Runner4 struct {
	ID              aping.RunnerID   `json:"selectionId"`
	Name            string           `json:"runnerName,omitempty"`
	Status          string           `json:"status,omitempty"` // ACTIVE, REMOVED, WINNER, LOSER, HIDDEN
	LastPriceTraded float64          `json:"lastPriceTraded,omitempty"`
	AvailableToBack aping.PriceSizes `json:"availableToBack,omitempty"`
	AvailableToLay  aping.PriceSizes `json:"availableToLay,omitempty"`
}


