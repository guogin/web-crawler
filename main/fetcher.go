package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

type Fetcher interface {
	Fetch(url string) FetchResult
}

type FetchResult struct {
	requestUrl   string
	responseBody string
	childUrls    []string
	error        error
}

type MyFetcher struct {
	httpClient *http.Client
}

func (f *MyFetcher) Fetch(url string) FetchResult {
	resp, err := f.httpClient.Get(url)

	if err != nil {
		return FetchResult{url, "", nil, fmt.Errorf("not found: %s", url)}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchResult{url, "", nil, fmt.Errorf("http status code: %v", resp.StatusCode)}
	}

	body, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(body)
	fmt.Println(bodyString)

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return FetchResult{url, "", nil, fmt.Errorf("error loading response body: %s", url)}
	}

	links := make([]string, 0)
	// Find all links and add them to links
	document.Find("a").Each(func (index int, element *goquery.Selection){
		// See if the href attribute exists on the element
		href, exists := element.Attr("href")
		if exists {
			links = append(links, normalize(href, url))
		}
	})

	return FetchResult{url, string(body), links, nil}
}

func normalize(href string, url string) string {
	var link string
	if strings.HasPrefix(href, "/") {
		if strings.HasSuffix(url, "/") {
			link = url + href
		} else {
			link = url + "/" + href
		}
	} else {
		link = href
	}

	return link
}
