package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"GoStudy/internal/config"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func user(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/user/")
	fmt.Fprintf(w, "Hi there, user %s!", name)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/", user)
	mux.HandleFunc("/", home)
	return mux
}

func main() {
	addr := config.Env("HTTP_ADDR", ":8080")
	log.Fatal(http.ListenAndServe(addr, newMux()))
}
