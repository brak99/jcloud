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
