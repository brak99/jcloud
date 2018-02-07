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

		fmt.Println("Request received: ", time.Now())

		const delay = 5000 * time.Millisecond

		time.Sleep(delay)

		password := r.URL.Query().Get("password")

		sha512 := sha512.New()

		sha512.Write([]byte(password))

		encoded := base64.URLEncoding.EncodeToString(sha512.Sum(nil))

		fmt.Printf("sha512:\t\t%s\n", encoded)

		fmt.Fprintf(w, encoded)
	}
}

func main() {

	handler := handlePassword()
	http.HandleFunc("/hash", handler)
	http.ListenAndServe(":8088", nil)
}
