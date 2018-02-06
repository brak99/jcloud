package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMonkey(t *testing.T) {

	passwords := make(chan string)

	handler := handlePassword(passwords)

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

	passwords := make(chan string)

	handler := handlePassword(passwords)

	r, _ := http.NewRequest("POST", "http://localhost:8088/?password=", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}

func TestMonkey_WrongEndpoint(t *testing.T) {

	passwords := make(chan string)

	handler := handlePassword(passwords)

	r, _ := http.NewRequest("POST", "http://localhost:8088/?nope=whatever", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

}
