package main

import (
	"fmt"
	"sync"
)

//
// Several solutions to the crawler exercise from the Go tutorial (https://tour.golang.org/concurrency/10)
//

//
// Serial crawler
//

func CrawlSerial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		CrawlSerial(u, fetcher, fetched)
	}
	return
}

//
// Concurrent crawler with Mutex and WaitGroup
//

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func (f *fetchState) CheckAndMark(url string) bool {
	//defer f.mu.Unlock()

	//f.mu.Lock()
	if f.fetched[url] {
		return true
	}
	f.fetched[url] = true
	return false
}

func mkFetchState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)
	return f
}

func CrawlConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	if f.CheckAndMark(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			CrawlConcurrentMutex(u, fetcher, f)
		}(u) // Without the u argument there is a race
	}
	done.Wait()
	return
}

//
// Concurrent crawler with channels
//

func dofetch(url1 string, ch chan []string, fetcher Fetcher) {
	body, urls, err := fetcher.Fetch(url1)
	if err != nil {
		fmt.Println(err)
		ch <- []string{}
	} else {
		fmt.Printf("found: %s %q\n", url1, body)
		ch <- urls
	}
}

func master(ch chan []string, fetcher Fetcher) {
	n := 1
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if _, ok := fetched[u]; ok == false {
				fetched[u] = true
				n += 1
				go dofetch(u, ch, fetcher)
			}
		}
		n -= 1
		if n == 0 {
			break
		}
	}
}

func CrawlConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url}
	}()
	master(ch, fetcher)
}

//
// main
//

func main() {
	fmt.Printf("=== Serial===\n")
	url := "http://golang.org/"
	CrawlSerial(url, fetcher, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	CrawlConcurrentMutex(url, fetcher, mkFetchState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	CrawlConcurrentChannel(url, fetcher)
}

//
// Fetcher
//

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
