# k8s-Cache
[![Go Report Card](https://goreportcard.com/badge/github.com/PedroMsTavares/k8s-cache)](https://goreportcard.com/report/github.com/PedroMsTavares/k8s-cache)

Cache in kubernetes using configmaps

# How does it work

Given an URL and config, it will fetch the url, save it in a configmap and update it every 2 minutes. That data can be acessed then by and http endpoint.

# Example

Your config has to be located in `config/services.yaml` .

```
Services:
- name: example # Name of the service
  url: https://reqres.in/api/users?page=2 # Url that will be fetched
  path: /myendpoint  # Path where k8s-cache will serve the result
  headers: # Any headers required for the request
    - key: "MyHeader"
      value: "Test"
```

In the example above, the app:

- fetch the data from `https://reqres.in/api/users?page=2` using the headers defined
- exposes the edpoint `http://localhost:8080/myendpoint` with the data
- refresh's the data every 2 minutes

## Notes

ConfigMaps only allow up to 1Mb of data. Considering this all the data is saved in binary to maximize the capacity.

By default the app will use the namespace `kube-system` , if you want to use your own namespace, set the environment variable `NAMESPACE` with the name of your namespace.

### TODO

- [ ] Add retry on failure
- [x] Docs
- [ ] Testing
