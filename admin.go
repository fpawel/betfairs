package main

import (
	"os"

)

var (
	adminBetfairUser = os.Getenv("BETFAIR_LOGIN_USER")
	adminBetfairPass = os.Getenv("BETFAIR_LOGIN_PASS")
)
