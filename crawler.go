package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if baseURL.Host != currentURL.Host {
		fmt.Println("different domains at first call!")
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, ok := pages[normCurrentURL]
	if ok {
		pages[normCurrentURL]++
		return
	} else {
		pages[normCurrentURL] = 1
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML %v\n", err)
		return
	}

	pageLinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Println("error getting urls from html")
		return
	}

	for _, link := range pageLinks {
		fmt.Printf("calling %s\n", link)
		crawlPage(rawBaseURL, link, pages)
	}

}
