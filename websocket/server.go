package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	go h.run()
	router.HandleFunc("/ws", myWs)
	addr := env("WEBSOCKET_ADDR", "127.0.0.1:8080")
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("err:", err)
	}
}
