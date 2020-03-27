# k8s-cache

The intent of this project is to create a proxy with cache taking advantage of kubernetes configmaps.

## How does it work ?

Given the yaml configuration bellow, it fetch`s the request and saves it in a configmap.

```
Services:
- name: example
  url: https://reqres.in/api/users?page=2
#  Add headers to your request
#  headers:
#    - key: "test"
#      value: "value "
```

Under the hood it updates the configmap every 5mins.

Exposes an http endpoint that give access to service/configmap.

`localhost:8080/service/{nameoftheservice}`
