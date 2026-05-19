package main

import (
	"fmt"
	"log"
	"net/http"

	"GoStudy/internal/config"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, user %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/", user)
	addr := config.Env("HTTP_ADDR", ":8080")
	log.Fatal(http.ListenAndServe(addr, nil))
}
