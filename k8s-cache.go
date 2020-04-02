package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigParse is used to parse the configuration from the yaml file
func ConfigParse() (s Services, err error) {

	services := Services{}
	data, err := ioutil.ReadFile("config/services.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(data), &services)
	if err != nil {
		fmt.Println(err)
	}
	return services, err
}

// ProcessConfig Process all the configs and create the configmaps
func ProcessConfig() error {
	fmt.Println("Start to processing")
	services, err := ConfigParse()
	if err != nil {
		fmt.Println(err)
		return err
	}
	httpClient := &http.Client{}

	// get service by service
	for _, service := range services.Services {
		req, err := http.NewRequest("GET", service.URL, nil)
		for _, header := range service.Headers {
			req.Header.Add(header.Key, header.Value)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if resp.StatusCode != 200 {
			fmt.Println("Failed to fetch service %s", service.Name)
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		k8sClient, err := ConfigK8s()
		if err != nil {
			fmt.Println(err)
			return err
		}
		namespace := os.Getenv("NAMESPACE")
		configmapClient := k8sClient.CoreV1().ConfigMaps(namespace)

		configmap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: service.Name,
			},
			BinaryData: map[string][]byte{
				service.Name: body,
			}}
		// need proper validatation an update
		_, err = configmapClient.Update(configmap)

		if err != nil {
			fmt.Println(err)
			//return err
		}
		_, err = configmapClient.Create(configmap)

		if err != nil {
			fmt.Println(err)
			return err
		}

	}
	fmt.Println(err)
	return nil
}

// SyncDaemon is responsenble to keep all the syncs working
func SyncDaemon() {
	for {
		go ProcessConfig()
		<-time.After(5 * time.Minute)
	}
}
