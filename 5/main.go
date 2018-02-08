package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

var idStore struct {
	ids       []string
	currentID int
	lock      sync.Mutex
}

type idJob struct {
	id       int
	password string
}

var workQueue = make(chan idJob, 100)

func handleShutdown(signal chan struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request received: %s", time.Now())
		fmt.Println("shutting down")

		signal <- struct{}{}
		fmt.Fprintf(w, "shutdown")
	}
}

func handlePassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		defer wg.Done()

		fmt.Printf("Request received: %s\n", time.Now())

		idStore.lock.Lock()

		id := idStore.currentID
		idStore.currentID++

		idStore.lock.Unlock()

		ret := fmt.Sprintf("%v", id)
		fmt.Fprintf(w, ret)

		password := r.URL.Query().Get("password")
		idJobReq := idJob{id: id, password: password}

		workQueue <- idJobReq

		wg.Add(1)

	}
}

func handleGetPassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//TODO: should check method
		id := strings.TrimPrefix(r.URL.Path, "/hash/")
		fmt.Printf("%v\n", id)

		idInt, _ := strconv.ParseInt(id, 10, 64)
		fmt.Fprintf(w, idStore.ids[idInt])
	}
}

func initIDStore() {
	idStore.currentID = 0
	idStore.lock = sync.Mutex{}
	idStore.ids = make([]string, 500)
}

func storePasswordHash() {
	select {
	case req := <-workQueue:
		const delay = 5000 * time.Millisecond

		time.Sleep(delay)

		sha512 := sha512.New()

		sha512.Write([]byte(req.password))

		encoded := base64.StdEncoding.EncodeToString(sha512.Sum(nil))

		fmt.Printf("adding: %v\n", req.id)
		fmt.Printf("sha512:\t\t%s\n", encoded)

		idStore.ids[req.id] = encoded
	}
}

func main() {

	initIDStore()

	stop := make(chan struct{})

	handler := handlePassword()
	shutdown := handleShutdown(stop)
	getPassword := handleGetPassword()

	http.HandleFunc("/hash", handler)
	http.HandleFunc("/shutdown", shutdown)
	http.HandleFunc("/hash/", getPassword)

	server := &http.Server{Addr: ":8088"}

	go func() {
		server.ListenAndServe()
	}()

	go func() {
		for {
			storePasswordHash()
		}
	}()

	<-stop
	wg.Wait()

	server.Shutdown(context.Background())
}
