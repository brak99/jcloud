package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func handlePassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			fmt.Println("Request received: ", time.Now())

			const delay = 5000 * time.Millisecond

			time.Sleep(delay)

			password := r.URL.Query().Get("password")

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
	handler := handlePassword()
	http.HandleFunc("/hash", handler)
	http.ListenAndServe(":8088", nil)
}
