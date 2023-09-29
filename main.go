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

type UserSecretResponse struct {
	Secret string `json:"Secret"`
}

func testPort(serverIP string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", serverIP, port)

	conn, err := net.Dial("tcp", address)
	if err == nil {
		conn.Close()

		pingURL := fmt.Sprintf("http://%s:%d/ping", serverIP, port)
		respPing, err := http.Get(pingURL)
		if err == nil {
			defer respPing.Body.Close()
			fmt.Printf("Port %d accessible - GET Response for /ping: %s\n", port, respPing.Status)
		}

		signupURL := fmt.Sprintf("http://%s:%d/signup", serverIP, port)
		body := []byte(`{"User": "Igor"}`)
		respSignup, err := http.Post(signupURL, "application/json", bytes.NewBuffer(body))
		if err == nil {
			defer respSignup.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /signup: %s\n", port, respSignup.Status)
		}

		checkURL := fmt.Sprintf("http://%s:%d/check", serverIP, port)
		respCheck, err := http.Post(checkURL, "application/json", bytes.NewBuffer(body))
		if err == nil {
			defer respCheck.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /check: %s\n", port, respCheck.Status)

			responseBody, err := ioutil.ReadAll(respCheck.Body)
			if err == nil {
				fmt.Printf("Contenu de la réponse de /check : %s\n", string(responseBody))
			} else {
				fmt.Printf("Erreur lors de la lecture de la réponse de /check : %v\n", err)
			}
		}

		userRequestBody := []byte(`{"User": "Igor"}`)
		userSecretURL := fmt.Sprintf("http://%s:%d/getUserSecret", serverIP, port)
		respUserSecret, err := http.Post(userSecretURL, "application/json", bytes.NewBuffer(userRequestBody))
		if err == nil {
			defer respUserSecret.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /getUserSecret: %s\n", port, respUserSecret.Status)

			userSecret, err := ioutil.ReadAll(respUserSecret.Body)
			if err == nil {
				fmt.Printf("Secret de l'utilisateur : %s\n", string(userSecret))
			} else {
				fmt.Printf("Erreur lors de la lecture de la réponse de /getUserSecret : %v\n", err)
			}
		}

		userLevelRequestBody := []byte(`{"User": "Igor", "Secret": "773079ad807a9694223d69ea5c9a05b0e98a74044a0e5d72ad0fcfcd0b72f20b"}`)

		userLevelURL := fmt.Sprintf("http://%s:%d/getUserLevel", serverIP, port)
		respUserLevel, err := http.Post(userLevelURL, "application/json", bytes.NewBuffer(userLevelRequestBody))
		if err == nil {
			defer respUserLevel.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /getUserLevel: %s\n", port, respUserLevel.Status)
			userLevel, err := ioutil.ReadAll(respUserLevel.Body)
			if err == nil {
				fmt.Printf("Niveau de l'utilisateur : %s\n", string(userLevel))
			} else {
				fmt.Printf("Erreur lors de la lecture de la réponse de /getUserLevel : %v\n", err)
			}
		}

		userPointsRequestBody := []byte(`{"User": "Igor", "Secret": "773079ad807a9694223d69ea5c9a05b0e98a74044a0e5d72ad0fcfcd0b72f20b"}`)
		userPointsURL := fmt.Sprintf("http://%s:%d/getUserPoints", serverIP, port)
		respUserPoints, err := http.Post(userPointsURL, "application/json", bytes.NewBuffer(userPointsRequestBody))
		if err == nil {
			defer respUserPoints.Body.Close()
			fmt.Printf("Port %d accessible - POST Response for /getUserPoints: %s\n", port, respUserPoints.Status)
			userPoints, err := ioutil.ReadAll(respUserPoints.Body)
			if err == nil {
				fmt.Printf("Points de l'utilisateur : %s\n", string(userPoints))
			} else {
				fmt.Printf("Erreur lors de la lecture de la réponse de /getUserPoints : %v\n", err)
			}
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
	wg.Wait()
}
