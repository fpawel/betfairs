package football4

import (
	"io"
	"encoding/binary"
	"runtime"
	"log"
	"path/filepath"
)


func (x Game) Serialize(w io.Writer)  {
	writeInt(w,x.ID)
}

func writeInt(w io.Writer, value int) {
	b := make([]byte,8)
	binary.PutVarint(b, int64(value))
	_,err := w.Write(b)
	check(err)
}

func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Panicf("%s:%d %v\n", filepath.Base(file), line, err)
	}
}
