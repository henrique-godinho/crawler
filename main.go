package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("# usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	if len(os.Args[1:]) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	if len(os.Args[1:]) == 3 {
		fmt.Printf("starting crawl of: %s\n", os.Args[1])
	}

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("error gettings maxConcurrency")
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("error getting maxPages")
		os.Exit(1)
	}

	pages := make(map[string]PageData)
	baseURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("error parsing baseURL in main: %v", err)
		os.Exit(1)
	}
	concurrencyControl := make(chan struct{}, maxConcurrency)

	cfg := config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: concurrencyControl,
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	cfg.wg.Add(1)
	go cfg.crawlPage(os.Args[1])
	cfg.wg.Wait()

	for page := range cfg.pages {
		fmt.Println(page)
	}

}
