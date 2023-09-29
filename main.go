package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

func testPort(serverIP string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", serverIP, port)

	// Tentative de connexion au serveur
	conn, err := net.Dial("tcp", address)
	if err == nil {
		conn.Close()

		// Faire une requête HTTP GET pour /ping
		pingURL := fmt.Sprintf("http://%s:%d/ping", serverIP, port)
		respPing, err := http.Get(pingURL)
		if err == nil {
			defer respPing.Body.Close()
			fmt.Printf("Port %d accessible - GET Response for /ping: %s\n", port, respPing.Status)
		}

		// Faire une requête HTTP GET pour /pong
		pongURL := fmt.Sprintf("http://%s:%d/pong", serverIP, port)
		respPong, err := http.Get(pongURL)
		if err == nil {
			defer respPong.Body.Close()
			fmt.Printf("Port %d accessible - GET Response for /pong: %s\n", port, respPong.Status)
		}
	}
}

func main() {
	serverIP := "10.49.122.144"
	minPort := 1
	maxPort := 10000

	var wg sync.WaitGroup

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		go testPort(serverIP, port, &wg)
	}

	// Attendre que toutes les goroutines se terminent
	wg.Wait()
}
