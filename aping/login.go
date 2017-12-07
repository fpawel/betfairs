package aping

import (
	"fmt"
	"net/http"
	"regexp"
)

// Login выполняет авторизацию на  betfair.com
func Login(user, pass string) (sessionToken string, err error) {

	const urlPattern = `https://identitysso.betfair.com/api/login?username=%s&password=%s&login=true&redirectMethod=POST&product=home.betfair.int&url=https://www.betfair.com/`

	urlStr := fmt.Sprintf(urlPattern, user, pass)
	var req *http.Request
	if req, err = http.NewRequest("POST", urlStr, nil); err != nil {
		return
	}

	var client http.Client
	var response *http.Response
	if response, err = client.Do(req); err != nil {
		return
	}
	strSetCookie := response.Header.Get("Set-Cookie")

	m := regexp.MustCompile("ssoid=([^;]+);").FindStringSubmatch(strSetCookie)
	if len(m) < 2 {
		err = fmt.Errorf("no headers in response %v", strSetCookie)
		return
	}
	sessionToken = m[1]
	return

}
