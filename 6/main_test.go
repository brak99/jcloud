package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMonkey(t *testing.T) {

	handler := handlePassword()

	r, _ := http.NewRequest("POST", "http://localhost:8088/?password=angryMonkey", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}

func TestMonkey_Empty(t *testing.T) {

	handler := handlePassword()

	r, _ := http.NewRequest("POST", "http://localhost:8088/?password=", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}

func TestShutdown(t *testing.T) {

	stop := make(chan struct{})

	handler := handleShutdown(stop)

	r, _ := http.NewRequest("POST", "http://localhost:8088/shutdown?password=angryMonkey", nil)

	w := httptest.NewRecorder()

	go func() {
		handler(w, r)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header.Get("Content-Type"))
		fmt.Println(string(body))

	}()

	<-stop

}

func TestGetPassword(t *testing.T) {

	initIDStore()

	idStore.currentID = 2
	idStore.ids[0] = "something"
	idStore.ids[1] = "bob"

	handler := handleGetPassword()

	r, _ := http.NewRequest("GET", "http://localhost:8088/hash/1", nil)

	w := httptest.NewRecorder()

	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	if string(body) != "bob" {
		t.Fatal("response should be 'bob'")
	}

}
func TestInit(t *testing.T) {

	initIDStore()

	if idStore.currentID != 0 {
		t.Fatal("currentID != 0")
	}

	if idStore.ids == nil {
		t.Fatal("ids not initialized")
	}
}

func TestStorePasswordHash(t *testing.T) {
	//setup
	initIDStore()

	idJobReq := idJob{id: 0, password: "angryMonkey"}

	workQueue <- idJobReq

	storePasswordHash()

	if idStore.ids[0] != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
		t.Fatal("stored password does not equal ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==")
	}
}

func TestStats404(t *testing.T) {

	handler := handleStatistics()

	r, _ := http.NewRequest("POST", "http://localhost:8088/stats", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 404 {
		t.Fatal("response should be 404, not %v", resp.StatusCode)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}

func TestInitStats(t *testing.T) {

	initStats()

	if stats.total != 0 {
		t.Fatal("total != 0")
	}

	if stats.totalRequestTime != 0 {
		t.Fatal("titalRequestTime != 0")
	}
}

func TestStats(t *testing.T) {

	handler := handleStatistics()

	stats.total = 5
	stats.totalRequestTime = 471

	r, _ := http.NewRequest("GET", "http://localhost:8088/stats", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatal("response should be 200, not ", resp.StatusCode)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}

func TestUpdateStats(t *testing.T) {
	stats.total = 5
	stats.totalRequestTime = 100

	updateStats(5000)

	if stats.totalRequestTime != 105 {
		t.Fatal("totalRequestTime should be 105")
	}
}
