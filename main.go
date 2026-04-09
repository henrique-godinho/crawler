package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args[1:]) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	if len(os.Args[1:]) == 1 {
		fmt.Printf("starting crawl of: %s\n", os.Args[1])
	}

	pages := make(map[string]int)
	crawlPage(os.Args[1], os.Args[1], pages)

	for page, count := range pages {
		fmt.Printf("page: %s -> count: %d\n", page, count)
	}

}
