package football2

import (
	"encoding/json"
	"hash/fnv"
	"log"
	"strconv"
)

// GamesChanges содержит данные об изменения в списке игр
type GamesChanges struct {
	Reset bool          `json:"reset,omitempty"`
	New   Games         `json:"new,omitempty"`
	Out   []int         `json:"out,omitempty"`
	Upd   []GameChanges `json:"upd,omitempty"`
}

func (x *GamesChanges) Empty() bool {
	return len(x.New) == 0 && len(x.Out) == 0 && len(x.Upd) == 0
}

func (x *GamesChanges) GetHashCode() string {
	fnv32a := fnv.New64a()
	bytes, err := json.Marshal(x)
	if err != nil {
		log.Fatal("json.Marshal GamesChanges")
	}
	fnv32a.Write(bytes)

	return strconv.FormatUint(fnv32a.Sum64(), 16)
}

func (x Games) Changes(next Games) (r GamesChanges) {
	if len(x) == 0{
		r.Reset = true
		r.New = next
		return
	}


	mPrev := make(map[int]Game)
	mNext := make(map[int]Game)

	for _, game := range x {
		mPrev[game.ID] = game
	}

	for _, game := range next {
		mNext[game.ID] = game

		if _, ok := mPrev[game.ID]; !ok {
			r.New = append(r.New, game)
		}
	}

	for id, gamePrev := range mPrev {
		gameNext, ok := mNext[id]
		if !ok {
			r.Out = append(r.Out, id)
		} else {
			changes := gamePrev.Changes(gameNext)
			if !changes.Empty() {
				r.Upd = append(r.Upd, changes)
			}
		}
	}
	return
}
