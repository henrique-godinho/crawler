package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ok := cfg.addPageVisit(normCurrentURL); ok {
		fmt.Printf("crawling %s\n", rawCurrentURL)
		html, err := getHTML(rawCurrentURL)
		if err != nil {
			fmt.Printf("error getting HTML %v\n", err)
			return
		}

		cfg.mu.Lock()
		cfg.pages[normCurrentURL] = extractPageData(html, rawCurrentURL)
		cfg.mu.Unlock()

		pageLinks, err := getURLsFromHTML(html, cfg.baseURL)
		if err != nil {
			fmt.Println("error getting urls from html")
			return
		}

		for _, link := range pageLinks {
			cfg.wg.Add(1)
			go cfg.crawlPage(link)
		}

	}

}

func (cfg *config) addPageVisit(normalisedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if len(cfg.pages) >= cfg.maxPages {
		return false
	}
	_, ok := cfg.pages[normalisedURL]
	if !ok {
		cfg.pages[normalisedURL] = PageData{}
		return true
	}
	return false
}
