package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ServiceHandler handler to serve http endpoints with data
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	k8sClient, err := ConfigK8s()
	if err != nil {
		log.Error(err)
	}
	namespace := GetNamespace()
	rcm, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(config[r.URL.Path], metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		contentType := r.Header.Get("Content-Type")
		if (len(contentType) == 0) {
			contentType = "application/json"
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		bodyString := string(rcm.BinaryData[config[r.URL.Path]])
		fmt.Fprintf(w, bodyString)
		log.WithFields(log.Fields{
			"Service": config[r.URL.Path],
			"Path":    r.URL.Path,
		}).Info("Request done")

	}
}

//ReadyHandler Readyness probe
func ReadyHandler(w http.ResponseWriter, r *http.Request) {

	if ready == true {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Found at least one ConfigMap")
		return
	}else {
		k8sClient, err := ConfigK8s()
		if err != nil {
			log.Error(err)
		}
		namespace := GetNamespace()

		services, err := ConfigParse()
		if err != nil {
			log.Fatal(err)
		}
		configmapClient := k8sClient.CoreV1().ConfigMaps(namespace)

		// get service by service
		for _, service := range services.Services {
			_, err = configmapClient.Get(service.Name, metav1.GetOptions{})
			if err == nil {

				ready = true
			}
		}
	}
	if ready == false {
		go ProcessConfig()
		w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Can't found ConfigMaps")
			return
	}else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Found at least one ConfigMap")
			return
	}
}

//RootHandler default route
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>K8s-cache</h1><a href='https://github.com/PedroMsTavares/k8s-cache'>https://github.com/PedroMsTavares/k8s-cache</a>")
}
