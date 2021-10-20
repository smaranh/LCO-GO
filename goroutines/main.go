package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func main() {
	websiteList := []string{
		"https://google.com",
		"https://fb.com",
		"https://lco.dev",
		"https://goat.com",
		"https://github.com",
	}
	for _, site := range websiteList {
		go getStatusCode(site)
		wg.Add(1)
	}

	wg.Wait()
}

func getStatusCode(endpoint string) {
	defer wg.Done()
	res, err := http.Get(endpoint)
	if err != nil {
		log.Fatal("Failed to reach the endpoint: ", err)
	}

	fmt.Printf("Status code is %d for endpoint %s\n", res.StatusCode, endpoint)
}
