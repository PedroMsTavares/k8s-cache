package main

import (
	"testing"
)

func TestConfigParse(t *testing.T) {
	svc, err := ConfigParse()

	if err != nil {
		t.Error(err)
	}

	if !(svc.Services[0].Name == "example" && svc.Services[0].URL == "https://reqres.in/api/users?page=2" && svc.Services[0].Path == "/example") {
		t.Error("Config values don't match the example")
	}
}

func TestGetrequest(t *testing.T) {

	s := &Service{
		Name: "Test",
		URL:  "https://httpbin.org/get",
		Path: "",
	}

	_, err := GetRequest(s)

	if err != nil {
		t.Error(err)
	}

	s.URL = "https://httpbin.org/status/404"

	_, err = GetRequest(s)

	if err == nil {
		t.Error("Expecting a log for a 404 here")
	}

}
