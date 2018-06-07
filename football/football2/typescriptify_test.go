package football2

import (
	"testing"
	"github.com/OneOfOne/struct2ts"
	"github.com/fpawel/betfairs/football"
	"os"
)

func TestTypescriptify(t *testing.T){

	f, err := os.Create (`F:\Frontend\betfairf\src\football2_.ts`)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	s2ts := struct2ts.New(&struct2ts.Options{Indent:"    ",})

	for _,x := range []interface{} {
		football.Game{},
		Game{},
		GameChanges{},
		GamesChanges{},
	}{
		s2ts.Add(x)
	}

	s2ts.RenderTo(f)


}
