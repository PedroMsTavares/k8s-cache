package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ConfigK8s K8s client
func ConfigK8s() (*kubernetes.Clientset, error) {

	incluster := os.Getenv("INCLUSTER")

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	if incluster == "FALSE" {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			panic(err.Error())
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}
