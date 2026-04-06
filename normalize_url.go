package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %w", err)
	}

	host := strings.ToLower(u.Host)
	path := strings.ToLower(strings.TrimRight(u.Path, "/"))

	normalizedURL := host + path
	return normalizedURL, nil
}
