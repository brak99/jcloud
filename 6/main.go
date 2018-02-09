package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
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
	id        int
	password  string
	timeSoFar int64
}

var stats struct {
	total            int64
	totalRequestTime int64
	lock             sync.Mutex
}

//StatResponse is the response returned for stats requests
type StatResponse struct {
	Total   int64
	Average float64
}

var workQueue = make(chan idJob, 5000)

func handleShutdown(signal chan struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			fmt.Println("Request received: %s", time.Now())
			fmt.Println("shutting down")

			signal <- struct{}{}
			fmt.Fprintf(w, "shutdown")
		} else {
			http.NotFound(w, r)
		}
	}
}

func updateStats(nanoseconds int64) {
	stats.lock.Lock()
	stats.total++
	stats.totalRequestTime += nanoseconds / int64(1000)
	stats.lock.Unlock()
}

func handlePassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			fmt.Printf("Request received: %s\n", time.Now())

			start := time.Now()

			idStore.lock.Lock()

			id := idStore.currentID
			idStore.currentID++

			idStore.lock.Unlock()

			ret := fmt.Sprintf("%v", id)
			fmt.Fprintf(w, ret)

			//add hash job to queue to be processed
			password := r.URL.Query().Get("password")
			elapsed := time.Since(start)
			idJobReq := idJob{id: id, password: password, timeSoFar: elapsed.Nanoseconds()}

			workQueue <- idJobReq

			wg.Add(1)

		} else {
			http.NotFound(w, r)
		}
	}
}

func handleGetPassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			fmt.Printf("Request received: %s\n", time.Now())

			//extract the id from the url
			id := strings.TrimPrefix(r.URL.Path, "/hash/")
			fmt.Printf("%v\n", id)

			idInt, _ := strconv.ParseInt(id, 10, 64)
			fmt.Fprintf(w, idStore.ids[idInt])

		} else {
			http.NotFound(w, r)
		}
	}
}

func handleStatistics() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			fmt.Printf("Request received: %s\n", time.Now())

			average := float64(stats.totalRequestTime) / float64(stats.total)

			resp := StatResponse{Total: stats.total, Average: average}

			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		} else {
			http.NotFound(w, r)
		}
	}
}

func initIDStore() {
	idStore.currentID = 0
	idStore.lock = sync.Mutex{}
	idStore.ids = make([]string, 500)
}

func initStats() {
	stats.total = 0
	stats.totalRequestTime = 0
}

func storePasswordHash() {
	select {
	case req := <-workQueue:
		const delay = 5000 * time.Millisecond

		time.Sleep(delay)

		start := time.Now() // starting here instead of including the 5s wait delay

		sha512 := sha512.New()

		sha512.Write([]byte(req.password))

		encoded := base64.StdEncoding.EncodeToString(sha512.Sum(nil))

		fmt.Printf("adding: %v\n", req.id)
		fmt.Printf("sha512:\t\t%s\n", encoded)

		idStore.ids[req.id] = encoded

		wg.Done()

		elapsed := time.Since(start)

		updateStats(req.timeSoFar + elapsed.Nanoseconds())
	}
}

func main() {

	initIDStore()
	initStats()

	stop := make(chan struct{})

	handler := handlePassword()
	shutdown := handleShutdown(stop)
	getPassword := handleGetPassword()
	statistics := handleStatistics()

	http.HandleFunc("/hash", handler)
	http.HandleFunc("/shutdown", shutdown)
	http.HandleFunc("/hash/", getPassword)
	http.HandleFunc("/stats", statistics)

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

	server.Shutdown(context.Background())

	//wait for all the hashes to be saved (if we had a backing store like Redis or something)
	wg.Wait()
}
