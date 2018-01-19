package main

import (
	"log"
	"path/filepath"
	"runtime"
)

func main (){
	daemon()
}


func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Panicf("%s:%d %v\n", filepath.Base(file), line, err)
	}
}