package countries

import (
	"testing"

	"fmt"
	"os"
	"bufio"
	"strings"
)

func TestFetch(t *testing.T){

	f,err := os.Open("countries.tab")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	fOut, err := os.Create("data.go")
	if err != nil {
		t.Fatal(err)
	}
	defer fOut.Close()

	fOut.WriteString(`
package countries

var countries [] Country = {

`)

	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line,"\t")
		if len(s) != 8{
			t.Fatal("BAD LINE:",line)
		} else {
			fmt.Fprintf(fOut,
				`
	{
		Name : %q,
		FullName : %q,
		English : %q,
		Alpha2 : %q,
		Alpha3 : %q,
		ISO : %q,
		Location : %q,
		LocationPrecise : %q,
	},`,
				s[0], s[1], s[2],
				s[3], s[4],
				s[5], s[6], s[7] )
		}
	}

	fOut.WriteString(`}`)

}


