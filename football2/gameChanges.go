package football2

// GameChanges содержит изменения в футбольной игре
type GameChanges struct {
	ID             int      `json:"id"`
	Order          *int     `json:"order,omitempty"`
	ScoreHome      *int     `json:"score_home,omitempty"`
	ScoreAway      *int     `json:"score_away,omitempty"`
	InPlay         *bool    `json:"in_play,omitempty"`
	Time           *string  `json:"time,omitempty"`
	Competition *string `json:"competition,omitempty"`
	Country *string `json:"country,omitempty"`
	WinBack        *float64 `json:"win_back,omitempty"`
	WinLay         *float64 `json:"win_lay,omitempty"`
	DrawBack       *float64 `json:"draw_lay,omitempty"`
	DrawLay        *float64 `json:"draw_back,omitempty"`
	LoseBack       *float64 `json:"lose_lay,omitempty"`
	LoseLay        *float64 `json:"lose_back,omitempty"`
	TotalMatched   *float64 `json:"total_matched,omitempty"`
	TotalAvailable *float64 `json:"total_available,omitempty"`
}

func (x GameChanges) Empty() bool {
	return x.Order == nil &&
		x.ScoreHome == nil &&
		x.ScoreAway == nil &&
		x.InPlay == nil &&
		x.Time == nil &&
		x.Country == nil &&
		x.Competition == nil &&
		x.WinBack == nil &&
		x.WinLay == nil &&
		x.LoseBack == nil &&
		x.LoseLay == nil &&
		x.DrawBack == nil &&
		x.DrawLay == nil &&
		x.TotalMatched == nil &&
		x.TotalAvailable == nil

}
