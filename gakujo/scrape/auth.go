package scrape

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

// return RelayState, SAMLResponse
func RelayStateAndSAMLResponse(r io.Reader) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", "", err
	}
	selection := doc.Find("html > body > form > div > input")
	relayState, ok := selection.Attr("value")
	if !ok {
		return "", "", &ErrNotFound{Name: "RelayState"}
	}
	selection = selection.Next()
	samlResponse, ok := selection.Attr("value")
	if !ok {
		return "", "", &ErrNotFound{Name: "SAMLResponse"}
	}

	return relayState, samlResponse, nil
}

func ApacheToken(r io.Reader) (string, error) {
	// ページによってtokenの場所が違う場合
	selectors := []string{
		"#SC_A01_06 > form:nth-child(15) > div > input[type=hidden]",
		"#header > form:nth-child(4) > div > input[type=hidden]",
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	for _, selector := range selectors {
		selection := doc.Find(selector)
		token, ok := selection.Attr("value")
		if ok {
			return token, nil
		}
	}
	return "", &ErrNotFound{Name: "org.apache.struts.taglib.html.TOKEN"}
}
