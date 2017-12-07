package webclient


import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"compress/gzip"
)

type ErrorParserHTML struct {
	Html string
	err error
}

func (x ErrorParserHTML) Error() string{
	return x.err.Error()
}

func Fetch( URLStr string, parse func(*goquery.Document) error ) error {

	request := NewHTTPRequest(URLStr)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return  err
	}
	defer response.Body.Close()

	reader := response.Body

	if response.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return  err
		}
		reader = gzipReader
		defer gzipReader.Close()
	}
	doc,err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return  err
	}
	err = parse(doc)
	if err != nil {
		str,_ := doc.Html()
		return ErrorParserHTML{str,err}
	}
	return err
}


