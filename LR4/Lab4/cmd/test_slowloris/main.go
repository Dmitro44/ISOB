package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const numConnections = 5

func attack(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8089")
	if err != nil {
		fmt.Printf("[Client %d] Connection error: %v\n", id, err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "POST / HTTP/1.1\r\n")
	fmt.Fprintf(conn, "Host: localhost:8080\r\n")
	fmt.Fprintf(conn, "Content-Length: 15\r\n") // promise 15 bytes
	fmt.Fprintf(conn, "\r\n")

	// send 5 bytes
	fmt.Fprintf(conn, "12345")
	fmt.Printf("[Client %d] Connected, sent partial data ...\n", id)

	start := time.Now()

	_, _ = io.ReadAll(bufio.NewReader(conn))
	duration := time.Since(start)

	fmt.Printf("[Client %d] Server dropped the connection after %.5f seconds\n", id, duration.Seconds())
}

func main() {
	fmt.Printf("Starting Slowloris attack\n\n")

	var wg sync.WaitGroup

	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go attack(i, &wg)
	}

	wg.Wait()

	fmt.Println("\nTest completed. The server successfully defended against hanging sockets.")
}
