package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Large Payload Test")

	largeData := make([]byte, 10*1024*1024)

	resp, err := http.Post("http://localhost:8089/readBody", "application/octet-stream", bytes.NewReader(largeData))
	if err != nil {
		fmt.Printf("Attack stopped by server. Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Server responded with status: %s\n", resp.Status)
}
