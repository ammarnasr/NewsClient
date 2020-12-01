package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	// A Tokenizer returns a stream of HTML Tokens
	db := createDBObject(username, password, dbName)
	db.Query("Delete from links")
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
					newsEntries := findTokensInsideDiv(tokenStream, "a", "href")
					for _, ne := range newsEntries {
						addNewsEntry(db, ne)
					}
					db.Close()
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
		for _, divAttr := range t.Attr { //search for the 'id' attribute and expected value
			if divAttr.Key == attrbiuteKey && divAttr.Val == attrbiutevalue {
				fmt.Println("Found Div with Id : ", divAttr.Val)
				return true
			}
		}
	}
	return false
}

func findTokensInsideDiv(ts *html.Tokenizer, tag string, attrbiuteKey string) []newsEntry {
	var newsEntries []newsEntry
	for {
		if ts.Next() == html.ErrorToken {
			return newsEntries //End of Document
		}
		t := ts.Token()
		if t.Data == "div" {
			return newsEntries //End of Div
		}
		if t.Data == tag {
			resetEntry := true
			ne := newsEntry{}
			for _, divAttr := range t.Attr {

				if resetEntry { //start new entry
					ne = newsEntry{}
				}

				if divAttr.Key == attrbiuteKey { // add href
					ne.hyperRef = divAttr.Val
					resetEntry = false
				}

				if divAttr.Key == "title" { // add title
					ne.title = divAttr.Val
					newsEntries = append(newsEntries, ne)
					resetEntry = true
				}
			}
		}
	}
}
