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
