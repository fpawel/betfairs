package aping

import (
	"time"
	"strconv"
)

type EventType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
}

type Side int

const (
	Back Side = iota
	Lay
)

type MarketID string
type RunnerID int

// The competition the market is contained within. Usually only applies to Football competitions
type Competition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Event struct {

	// The unique id for the event
	ID string `json:"id"`

	// The name of the event
	Name string `json:"name"`

	// The ISO-2 code for the event.
	// A list of ISO-2 codes is available via
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	CountryCode string `json:"countryCode,omitempty"`

	// This is timezone in which the event is taking place./
	Timezone string `json:"timezone,omitempty"`

	Venue string `json:"venue,omitempty"`

	// The scheduled start date and time of the event.
	// This is Europe/London (GMT) by default
	OpenDate time.Time `json:"openDate"`
}

type MarketFilter struct {
	EventIDs []int `json:"eventIds"`
}

func (x MarketID) Int() int {
	n, _ :=  strconv.Atoi( string(x[2:]) )
	return n
}