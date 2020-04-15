package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)
	k8sClient, err := ConfigK8s()
	if err != nil {
		log.Error(err)
	}
	namespace := os.Getenv("NAMESPACE")
	rcm, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(Config[r.URL.Path], metav1.GetOptions{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		bodyString := string(rcm.BinaryData[Config[r.URL.Path]])
		fmt.Fprintf(w, bodyString)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello I am a cache for you !")
}
