package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ianbibby/sitemap"
)

func main() {
	urlFlag := flag.String("url", "", "URL to build a sitemap from")
	depth := flag.Int("depth", 1, "Max depth for the sitemap")
	flag.Parse()

	if strings.TrimSpace(*urlFlag) == "" {
		fmt.Fprintln(os.Stderr, "url is required")
		os.Exit(1)
	}

	s, err := sitemap.Build(*urlFlag, *depth)
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}
