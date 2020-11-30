package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type tokenInfo struct {
	tokenData string
}

func main() {
	// A Tokenizer returns a stream of HTML Tokens
	tokenStream := getHTMLTokenStreamFromURL("https://en.wikinews.org/wiki/Main_Page")

	for {

		tokenType := tokenStream.Next()

		switch tokenType {

		case html.ErrorToken:
			{
				// End of the document
				return
			}

		case html.StartTagToken: // Opening Tag
			{
				token := tokenStream.Token()
				tokenFound := findTokenByAttributes(token, "div", "id", "MainPage_latest_news_text")
				if tokenFound {
					findTokensInsideDiv(tokenStream, "a", "href")
					return
				}
			}
		}
	}
}

func getHTMLTokenStreamFromURL(url string) *html.Tokenizer {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return html.NewTokenizer(resp.Body)
}

func findTokenByAttributes(t html.Token, tag string, attrbiuteKey string, attrbiutevalue string) bool {
	if t.Data == tag {
		for _, divAttr := range t.Attr {
			if divAttr.Key == attrbiuteKey && divAttr.Val == attrbiutevalue {
				fmt.Println("Found Div with Id : ", divAttr.Val)
				return true
			}
		}
	}
	return false
}

func findTokensInsideDiv(ts *html.Tokenizer, tag string, attrbiuteKey string) {
	for {
		if ts.Next() == html.ErrorToken {
			return
		}
		t := ts.Token()
		if t.Data == "div" {
			return //End of Div
		}
		if t.Data == tag {
			for _, divAttr := range t.Attr {
				if divAttr.Key == attrbiuteKey {
					fmt.Println("Found Href : ", divAttr.Val)
				}
			}
		}
	}
}
