package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ConfigK8s K8s client
func ConfigK8s() (*kubernetes.Clientset, error) {

	incluster := os.Getenv("INCLUSTER")

	var config *rest.Config
	var err error

	if incluster == "" {
		config, err = rest.InClusterConfig()

		if err != nil {
			log.Panic(err)
		}
	}
	if incluster == "FALSE" {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			log.Panic(err)
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

//GetNamespace returns the value of the environment variable NAMESPACE or if not set the default the namespace
func GetNamespace() string {

	ns := os.Getenv("NAMESPACE")

	if ns == "" {
		ns = "kube-system"
	}
	return ns
}
