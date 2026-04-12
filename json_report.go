package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	var keys []string
	var sortedPages []PageData
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		page := pages[key]
		sortedPages = append(sortedPages, page)
	}

	data, err := json.MarshalIndent(sortedPages, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil

}
