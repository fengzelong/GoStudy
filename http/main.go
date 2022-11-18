package main

import (
	"fmt"
	"log"
	"net/http"
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
