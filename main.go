package main

import (
	"log"
	"path/filepath"
	"runtime"
	"os"
)

var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
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