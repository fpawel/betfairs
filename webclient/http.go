package webclient


import (
	"net/http"
	"log"
	"net/url"
)


func NewURL(URLpathStr string ) string {
	URLbase, err := url.Parse("https://www.betfair.com")
	if err != nil {
		log.Fatal(err)
	}
	URLpath, err := url.Parse(URLpathStr)
	if err != nil {
		log.Fatal(err)
	}
	return  URLbase.ResolveReference(URLpath).String()
}

func NewHTTPRequest(URLStr string) *http.Request {


	request, err := http.NewRequest("GET", URLStr, nil)
	if err != nil {
		log.Fatal(err)
	}

	const (
		ruCoockie  = `vid=39ad9e9d-12e6-487e-8f0c-fb89d881c015; bucket=2~0~test_search; wsid=13f06991-5f1a-11e6-862c-90e2ba0fa6a0; betexPtk=betexCurrency%3DGBP%7EbetexTimeZone%3DEurope%2FLondon%7EbetexRegion%3DGBR%7EbetexLocale%3Dru; mEWJSESSIONID=AA42CA3984085031B4C9F344A940BACB; betexPtkSess=betexCurrencySessionCookie%3DGBP%7EbetexRegionSessionCookie%3DGBR%7EbetexTimeZoneSessionCookie%3DEurope%2FLondon%7EbetexLocaleSessionCookie%3Dru%7EbetexSkin%3Dstandard%7EbetexBrand%3Dbetfair; PI=61999; pi=partner61999; UI=0; spi=0; bfsd=ts=1470847622600|st=p; _qst_s=1; _qsst_s=1470847831600; betfairSSC=lsSSC%3D1%3Bcookie-policy%3D1; _ga=GA1.2.783247160.1470847628; _gat=1; _qubitTracker=1470847628500.818360; _qubitTracker_s=1470847628500.818360; _qPageNum_betfair=1; _qst=%5B1%2C0%5D; _qsst=1470847840800; qb_ss_status=BOA5:Ma&OsT|BOBI:Ik&OsY|BOBQ:D0&Osa|BOBk:OP&Osd|BOKN:J&Ot6; _qb_se=BOA5:OsT&VZ1XPKE|BOBI:OsY&VZ1XPKE|BOBQ:Osa&VZ1WbWc|BOBk:Osd&VZ1XPKE|BOKN:Ot6&VZ1XPKE; qb_permanent=:0:0:0:0:0::0:1:0::::::::::::::::::::K6M&OMN&OsY&Osd&OuS&OsT&Osa&Ot6:VZ1XPKE; _q_geo=%5B%222%22%2C%2293.115.95.202%22%2C%22RO%22%2C%2212072%22%2C%22unknown%22%2C%2217843%22%2C%2244.4599%22%2C%2226.1333%22%5D; qb_cc=RO; update-browser=Wed%20Aug%2010%202016%2016%3A47%3A39%20GMT%2B0000%20(UTC); exp=ex; pref_md_pers_0="{\"com-es-info\":{\"spainRedirectNotification\":\"false\"}}"; ss_opts=BOA5:C&C|BOBI:C&C|BOBQ:B&B|BOBk:C&C|BOKN:C&C|_g:VZ1WbU4&VZ1XPIg&B&C`
	)


	httpHeader := http.Header{
		"Accept-Language": {"ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"},
		"Accept-Encoding": {"gzip,deflate,sdch"},

		"User-Agent": {"Mozilla/5.0 (Windows NT 6.1; rv:45.0) Gecko/20100101 Firefox/45.0"},

		"Accept":  {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Cookie":  {ruCoockie},
		"Referer": {"http://www.betfair.com/ru/"},
	}

	for k, v := range httpHeader {
		request.Header[k] = v
	}
	return request

}


/*
func NewTorHttpClient() *http.Client{
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9150")
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to obtain proxy dialer: %v", err)
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	return &http.Client{Transport: tbTransport}
}
*/


