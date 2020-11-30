package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	resp, _ := http.Get("https://en.wikinews.org/wiki/Main_Page")
	//bytes, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("HTML:\n\n", string(bytes))

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "div"
			if isAnchor {
				fmt.Println("We found a div!")

				for _, a := range t.Attr {
					if a.Key == "id" {
						if a.Val == "MainPage_latest_news_text" {
							fmt.Println("Found id:", a.Val)
						}
						break
					}
				}
			}
		}
	}
}
