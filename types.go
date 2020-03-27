package main

// Service struct that defines a service
type Service struct {
	Name    string    `yaml:"name"`
	URL     string    `yaml:"url"`
	Path    string    `yaml:"path"`
	Headers []*Header `yaml:"headers"`
}

// Header struct that defines the headers of a service
type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Services struct that contains all the services
type Services struct {
	Services []*Service `yaml:"Services"`
}
