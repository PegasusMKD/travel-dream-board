package main

import (
	"fmt"
	"net/url"
	"strings"
)

func getHostFromURL(URL string) string {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return ""
	}
	host := parsedURL.Hostname()
	host = strings.TrimPrefix(host, "www.")
	return host
}

func main() {
	fmt.Println(getHostFromURL("https://www.google.com/travel/flights"))
	fmt.Println(getHostFromURL("https://www.austrian.com/"))
	fmt.Println(getHostFromURL("https://www.esky.com/"))
}
