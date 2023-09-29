package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type CheckResponse struct {
	User string `json:"User"`
}

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

		// Faire une requête HTTP POST pour /signup avec le prénom
		signupURL := fmt.Sprintf("http://%s:%d/signup", serverIP, port)
		body := []byte(`{"User": "Igor"}`)
		respSignup, err := http.Post(signupURL, "application/json", bytes.NewBuffer(body))
		if err == nil {
			defer respSignup.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /signup: %s\n", port, respSignup.Status)
		}

		// Faire une requête HTTP POST pour /check avec le même corps JSON
		checkURL := fmt.Sprintf("http://%s:%d/check", serverIP, port)
		respCheck, err := http.Post(checkURL, "application/json", bytes.NewBuffer(body))
		if err == nil {
			defer respCheck.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /check: %s\n", port, respCheck.Status)

			// Lire le contenu de la réponse de /check
			responseBody, err := ioutil.ReadAll(respCheck.Body)
			if err == nil {
				fmt.Printf("Contenu de la réponse de /check : %s\n", string(responseBody))
			} else {
				fmt.Printf("Erreur lors de la lecture de la réponse de /check : %v\n", err)
			}
		}
	} else {
		fmt.Printf("Port %d inaccessible\n", port)
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
