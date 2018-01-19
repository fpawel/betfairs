package aping

///The static data about runners in a market
type RunnerCatalogue struct {

	// The unique id of the runner (selection)
	ID RunnerID `json:"selectionId"`

	// The name of the runner
	Name string `json:"runnerName,omitempty"`

	//The sort priority of this runner
	//SortPriority int `json:"sortPriority,omitempty"`

	// The handicap.  Enter the specific handicap value (returned by RUNNER in listMaketBook)
	// if the market is an Asian handicap market.
	//Handicap float64 `json:"handicap,omitempty"`
}
