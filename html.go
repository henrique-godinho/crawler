package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

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
				urls = append(urls, baseURL.Scheme+"://"+baseURL.Host+val)
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
				urls = append(urls, strings.TrimRight(baseURL.String(), "/")+val)
			}
		}
	})

	return urls, nil
}

func extractPageData(html, pageURL string) PageData {
	var pd PageData
	heading := getHeadingFromHTML(html)
	fp := getFirstParagraphFromHTML(html)
	pageURLasURL, err := url.Parse(pageURL)
	if err != nil {
		fmt.Println(err)
		return PageData{
			URL:            pageURL,
			Heading:        heading,
			FirstParagraph: fp,
			OutgoingLinks:  nil,
			ImageURLs:      nil,
		}
	}
	pd.URL = pageURL

	outLinks, err := getURLsFromHTML(html, pageURLasURL)
	if err != nil {
		fmt.Println(err)
		return PageData{}
	}
	imgURLS, err := getImagesFromHTML(html, pageURLasURL)
	if err != nil {
		fmt.Println(err)
		return PageData{}
	}
	pd.Heading = heading
	pd.FirstParagraph = fp
	pd.OutgoingLinks = outLinks
	pd.ImageURLs = imgURLS

	return pd
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		fmt.Println("fail to create req")
		return "", err
	}
	req.Header.Add("User-Agent", "BootCrawler/1.0")
	var client http.Client
	res, err := client.Do(req)

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("http error: %v", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("invalid content-type %s", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body %v", err)
	}

	return string(body), nil

}
