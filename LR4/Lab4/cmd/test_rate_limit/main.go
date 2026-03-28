package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Rate Limit Test")

	for i := 1; i <= 10; i++ {
		resp, err := http.Get("http://localhost:8089/")
		if err != nil {
			fmt.Printf("Response %d: Error: %v\n", i, err)
			continue
		}

		fmt.Printf("Response %d: Status code %d (%s) \n", i, resp.StatusCode, resp.Status)
		resp.Body.Close()
	}
}
