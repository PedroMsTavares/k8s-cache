package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

}

var Config map[string]string

func main() {

	r := mux.NewRouter()
	//r.HandleFunc("/service/{servicename}", ServiceHandler)

	// Register the diferent paths

	config, err := ConfigParse()
	if err != nil {
		log.Fatal("Can't parse config")
	}
	Config = make(map[string]string)
	for _, service := range config.Services {

		r.HandleFunc(service.Path, ServiceHandler)
		fmt.Println(service.Path)

		Config[service.Path] = service.Name
	}
	fmt.Println(Config)

	r.HandleFunc("/", RootHandler)

	go SyncDaemon()
	http.ListenAndServe(":8080", r)

}
