all: push

BUILDTAGS=

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
APP=cicd
USERSPACE?=k8s-community
RELEASE?=0.4.4
PROJECT?=github.com/${USERSPACE}/${APP}
GOOS?=linux

# App configuration: what service to listen to etc
SERVICE_HOST?=0.0.0.0
SERVICE_PORT?=8080
GITHUBINT_BASE_URL?=https://services.k8s.community/k8s-community/github-integration

REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
  COMMIT := git-$(shell git rev-parse --short HEAD)
endif

vendor: clean
	go get -u github.com/Masterminds/glide \
	&& glide install

build: vendor
	cd service \
	&& CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/version.RELEASE=${RELEASE} -X ${PROJECT}/version.COMMIT=${COMMIT} -X ${PROJECT}/version.REPO=${REPO_INFO}" \
		-o ../${APP}

install: build
	sudo ./${APP} install --service-host ${SERVICE_HOST} --service-port ${SERVICE_PORT} \
	 --githubint-base-url ${GITHUBINT_BASE_URL}

remove:
	sudo ./${APP} remove

stop:
	sudo ./${APP} stop

start:
	sudo ./${APP} start

fmt:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"gofmt -s -l {{.Dir}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

lint: utils
	@echo "+ $@"
	go get -u github.com/golang/lint/golint
	@go list -f '{{if len .TestGoFiles}}"golint {{.Dir}}/..."{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

vet:
	@echo "+ $@"
	@go vet $(shell go list ${PROJECT}/... | grep -v vendor)

test: vendor utils fmt lint vet
	@echo "+ $@"
	@go test -v -race -tags "$(BUILDTAGS) cgo" $(shell go list ${PROJECT}/... | grep -v vendor)

cover:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

clean:
	rm -f ${APP}
