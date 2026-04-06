package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	h1 := doc.Find("h1").Text()
	if len(h1) == 0 {
		h2 := doc.Find("h2").Text()
		return h2
	}
	return h1
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	main := doc.Find("main").ChildrenFiltered("p")
	if len(main.Text()) == 0 {
		return doc.Find("p").Text()
	}
	return main.Text()
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	urls := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		if val, ok := s.Attr("href"); ok {
			u, err := url.Parse(val)
			if err != nil {
				fmt.Println(err)
				return
			}
			if u.IsAbs() {
				urls = append(urls, u.String())
			} else {
				urls = append(urls, baseURL.String()+val)
			}
		}
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	urls := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		if val, ok := s.Attr("src"); ok {
			u, err := url.Parse(val)
			if err != nil {
				fmt.Println(err)
				return
			}
			if u.IsAbs() {
				urls = append(urls, u.String())
			} else {
				urls = append(urls, baseURL.String()+val)
			}
		}
	})

	return urls, nil
}
