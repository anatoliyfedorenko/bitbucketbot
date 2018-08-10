package main

import (
	"fmt"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func push(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Someone just pushed to repo!")
}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/push", push)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
