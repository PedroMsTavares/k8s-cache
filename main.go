package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	go SyncDaemon()
	r := mux.NewRouter()
	r.HandleFunc("/service/{servicename}", ServiceHandler)
	r.HandleFunc("/", RootHandler)
	http.ListenAndServe(":8080", r)

}
