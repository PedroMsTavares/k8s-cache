package main

type Service struct {
	Name    string    `yaml:"name"`
	URL     string    `yaml:"url"`
	Path    string    `yaml:"path"`
	Headers []*Header `yaml:"headers"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Services struct {
	Services []*Service `yaml:"Services"`
}
