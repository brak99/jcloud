package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func handleShutdown(signal chan struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			fmt.Println("Request received: %s", time.Now())
			fmt.Println("shutting down")

			signal <- struct{}{}
			fmt.Fprintf(w, "shutdown")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}

func handlePassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			defer wg.Done()

			fmt.Printf("Request received: %s", time.Now())

			password := r.URL.Query().Get("password")

			wg.Add(1)

			const delay = 5000 * time.Millisecond

			time.Sleep(delay)

			sha512 := sha512.New()

			sha512.Write([]byte(password))

			encoded := base64.StdEncoding.EncodeToString(sha512.Sum(nil))

			fmt.Printf("sha512:\t\t%s\n", encoded)

			fmt.Fprintf(w, encoded)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}

func main() {

	fmt.Println("Server started on port 8088")

	stop := make(chan struct{})

	handler := handlePassword()
	shutdown := handleShutdown(stop)

	http.HandleFunc("/hash", handler)
	http.HandleFunc("/shutdown", shutdown)

	server := &http.Server{Addr: ":8088"}

	go func() {
		server.ListenAndServe()
	}()

	<-stop
	wg.Wait()

	server.Shutdown(context.Background())
}
