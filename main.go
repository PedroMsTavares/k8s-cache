package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

// ConfigK8s K8s client
func ConfigK8s() (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	/*	// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", "/Users/p.tavares/.kube/config")
		if err != nil {
			panic(err.Error())
		} */

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
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
		configmapClient := k8sClient.CoreV1().ConfigMaps("k8s-cache")

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

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	k8sClient, err := ConfigK8s()
	if err != nil {
		fmt.Println(err)
	}

	rcm, err := k8sClient.CoreV1().ConfigMaps("k8s-cache").Get(vars["servicename"], metav1.GetOptions{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		bodyString := string(rcm.BinaryData[vars["servicename"]])
		fmt.Fprintf(w, bodyString)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello I am a cache for you !")
}

func main() {
	go SyncDaemon()
	r := mux.NewRouter()
	r.HandleFunc("/service/{servicename}", ServiceHandler)
	r.HandleFunc("/", RootHandler)
	http.ListenAndServe(":8080", r)

}
