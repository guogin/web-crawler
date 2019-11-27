package main

type Crawler struct {
	fetcher Fetcher
}

func (crawler *Crawler) Crawl(url string, depth int, callback func(result FetchResult)) {
	visited := SafeCounter{m: make(map[string]int)}
	fetcher := crawler.fetcher
	doRecursiveFetch(fetcher, url, depth, visited, callback)
	return
}

func doRecursiveFetch(fetcher Fetcher, url string, depth int, visited SafeCounter, callback func(result FetchResult)) {
	numberOfVisit := visited.GetAndIncrement(url)

	if numberOfVisit > 0 || depth <= 0 {
		return
	}

	result := fetcher.Fetch(url)
	callback(result)

	subtreeDone := make(chan bool, len(result.childUrls))
	for _, u := range result.childUrls {
		go func(url string) {
			doRecursiveFetch(fetcher, url, depth-1, visited, callback)
			subtreeDone <- true
		}(u)
	}

	// Wait for all subtree done
	for i := 0; i < len(result.childUrls); i++ {
		<-subtreeDone
	}

	return
}
