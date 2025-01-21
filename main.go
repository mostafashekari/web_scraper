package main

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func fetchTitleAndMeta(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Error parsing HTML of %s: %v", url, err)
		return
	}

	var title, meta string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "description" {
					for _, a := range n.Attr {
						if a.Key == "content" {
							meta = a.Val
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	results <- fmt.Sprintf("URL: %s\nTitle: %s\nMeta Description: %s\n", url, title, meta)
}

func main() {
	urls := []string{
		"https://golang.org",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.hobb.pro",
	}

	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go fetchTitleAndMeta(url, &wg, results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}
}
