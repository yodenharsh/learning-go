package main

import (
	"fmt"
	"net/url"
)

func main() {
	rawUrl := "http://example.com:8080/path?query=param#fragment"

	parsedUrl, err := url.Parse(rawUrl)

	if err != nil {
		fmt.Println("Error parsing URL")
	}

	fmt.Println("Parsed URL scheme: ", parsedUrl.Scheme)
	fmt.Println("Host: ", parsedUrl.Host)
	fmt.Println("Fragment: ", parsedUrl.Fragment)
	fmt.Println("Raw query: ", parsedUrl.RawQuery)

	rawUrl1 := "https://example.com/path?name=john&age=30"

	parsedUrl1, err := url.Parse(rawUrl1)

	if err != nil {
		fmt.Println("Could not parse parsedUrl1")
	}

	queryParams := parsedUrl1.Query()
	fmt.Println(queryParams)
	fmt.Println("Name: ", queryParams.Get("name"))
	fmt.Println("Age: ", queryParams.Get("age"))

	// Building a URL
	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/path",
	}

	query := baseUrl.Query()
	query.Set("name", "lona%21&")
	baseUrl.RawQuery = query.Encode()

	fmt.Println("Built URL: ", baseUrl.String())

	values := url.Values{}
	// Adding key-value pairs to the values object

	values.Add("name", "holy")
	values.Add("age", "16")
	values.Add("city", "mumbai")
	values.Add("country", "UK")

	encodedQuery := values.Encode()
	baseUrl1 := &url.URL{
		Scheme:   "https",
		Host:     "example.com",
		Path:     "/path",
		RawQuery: encodedQuery,
	}

	fmt.Println(baseUrl1)
}
