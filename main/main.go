package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: transport}

	myFetcher := &MyFetcher{httpClient: client}

	crawler := Crawler{fetcher: myFetcher}

	crawler.Crawl(args[0], 10, func(result FetchResult) {
		if result.error != nil {
			fmt.Println(result.error)
			return
		}
		fmt.Printf("Found: %s %q\n", result.requestUrl, result.responseBody)
	})

	return
}
