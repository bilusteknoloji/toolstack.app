package main

import (
	"log"
	"net/http"
	"os"
)

func getListenAddr() string {
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		return addr
	}

	return ":8000"
}

func main() {
	addr := getListenAddr()

	fs := http.FileServer(http.Dir("./site"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("Listening at", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
