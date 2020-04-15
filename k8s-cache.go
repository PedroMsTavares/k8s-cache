package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigParse is used to parse the configuration from the yaml file
func ConfigParse() (s Services, err error) {

	services := Services{}
	data, err := ioutil.ReadFile("config/services.yaml")
	if err != nil {
		log.Fatal("Can't read config file")
	}
	err = yaml.Unmarshal([]byte(data), &services)
	if err != nil {
		log.Error(err)
	}
	return services, err
}

func GetRequest(service *Service) (body []byte, err error) {

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", service.URL, nil)
	for _, header := range service.Headers {
		req.Header.Add(header.Key, header.Value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"URL":          service.URL,
			"ResponseCode": resp.StatusCode,
		}).Errorf("Failed to fetch service %s", service.Name)

		return nil, fmt.Errorf("Failed to fetch service %s", service.Name)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return body, err

}

// ProcessConfig Process all the configs and create the configmaps
func ProcessConfig() error {
	log.Info("Start to process services")
	services, err := ConfigParse()
	if err != nil {
		log.Error(err)
		return err
	}

	// get service by service
	for _, service := range services.Services {
		log.Infof("Start service %s sync...", service.Name)
		body, err := GetRequest(service)

		if err != nil {
			continue
		}
		k8sClient, err := ConfigK8s()
		if err != nil {
			log.Error(err)
			continue
		}

		// get the namespage where the cm will be created
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
			log.Info(err)
			_, err = configmapClient.Create(configmap)
		}

		if err != nil {
			log.Info(err)
			//return err
		}
		log.Infof("Service %s synced", service.Name)
	}
	return nil
}

// SyncDaemon is responsenble to keep all the syncs working
func SyncDaemon() {
	for {
		go ProcessConfig()
		<-time.After(2 * time.Minute)
	}
}
