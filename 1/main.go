package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		password := os.Args[1]

		sha512 := sha512.New()

		sha512.Write([]byte(password))

		encoded := base64.StdEncoding.EncodeToString(sha512.Sum(nil))

		fmt.Printf(encoded)
	} else {
		fmt.Printf("Takes a password as input and returns it after it has been hashed with SHA512.\n")
		fmt.Printf("Usage: go run main.go <password>\n")
	}

}
