[![Build Status](https://travis-ci.org/k8s-community/cicd.svg?branch=master)](https://travis-ci.org/k8s-community/cicd)

# cicd

The simplest CI/CD service to `make test` and `make build`.

## How to run

Only to try!

    env SERVICE_HOST=0.0.0.0 SERVICE_PORT=8080 go run ./service/cicd.go
    
    
Run tests:

    go test -v -race  $(go list ./... | grep -v vendor)