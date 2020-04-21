package main

import (
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

var config map[string]string
var ready bool

func main() {

	//readiness probe false by default
	ready = false

	// Parse /config/service.yaml
	c, err := ConfigParse()
	if err != nil {
		log.Fatal("Can't parse config")
	}

	// Create mux router
	r := mux.NewRouter()

	config = make(map[string]string)

	// Register all the custom roles to the services
	for _, service := range c.Services {

		r.HandleFunc(service.Path, ServiceHandler)

		config[service.Path] = service.Name
	}

	// Register default route
	r.HandleFunc("/", RootHandler)
	// Register ready route
	r.HandleFunc("/ready", ReadyHandler)

	// Start the sync daemon
	go SyncDaemon()

	// Start http server
	http.ListenAndServe(":8080", r)

}
