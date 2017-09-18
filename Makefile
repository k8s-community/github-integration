all: push

BUILDTAGS=

APP?=github-integration
PROJECT?=github.com/k8s-community/${APP}
REGISTRY?=registry.k8s.community
CA_DIR?=certs

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE?=0.3.4
GOOS?=linux
GOARCH?=amd64

GITHUBINT_LOCAL_PORT?=8080

NAMESPACE?=k8s-community
INFRASTRUCTURE?=stable
VALUES?=values-${INFRASTRUCTURE}

CONTAINER_IMAGE?=${REGISTRY}/${NAMESPACE}/${APP}
CONTAINER_NAME?=${APP}-${NAMESPACE}

REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
	COMMIT := git-$(shell git rev-parse --short HEAD)
endif

BUILDTAGS=

.PHONY: all
all: build

.PHONY: vendor
vendor: clean bootstrap
	dep ensure

.PHONY: build
build: vendor
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/version.RELEASE=${RELEASE} -X ${PROJECT}/version.COMMIT=${COMMIT} -X ${PROJECT}/version.REPO=${REPO_INFO}" \
		-o ${APP}

.PHONY: container
container: build
	@echo "+ $@"
	@docker build --pull -t $(CONTAINER_IMAGE):$(RELEASE) .

.PHONY: push
push: container
	@echo "+ $@"
	@docker push $(CONTAINER_IMAGE):$(RELEASE)

.PHONY: certs
certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name ${CONTAINER_NAME}-certs -d alpine:edge sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait ${CONTAINER_NAME}-certs
	@docker cp ${CONTAINER_NAME}-certs:/etc/ssl/certs/ca-certificates.crt ${CA_DIR}
	@docker rm -f ${CONTAINER_NAME}-certs
endif

run: container
	@echo "+ $@"
	@docker run --name ${CONTAINER_NAME} -p ${GITHUBINT_LOCAL_PORT}:${GITHUBINT_LOCAL_PORT} \
		-e "GITHUBINT_LOCAL_PORT=${GITHUBINT_LOCAL_PORT}" \
		-d $(CONTAINER_IMAGE):$(RELEASE)
	@sleep 1
	@docker logs ${CONTAINER_NAME}

.PHONY: deploy
deploy: push
	helm upgrade ${CONTAINER_NAME} -f charts/${VALUES}.yaml charts \
		--kube-context ${INFRASTRUCTURE} --namespace ${NAMESPACE} \
		--version=${RELEASE} -i --wait

GO_LIST_FILES=$(shell go list ${PROJECT}/... | grep -v vendor)

.PHONY: fmt
fmt:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"gofmt -s -l {{.Dir}}"{{end}}' ${GO_LIST_FILES} | xargs -L 1 sh -c

.PHONY: lint
lint: bootstrap
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"golint {{.Dir}}/..."{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet $(shell go list ${PROJECT}/... | grep -v vendor)

.PHONY: test
test: vendor fmt lint vet
	@echo "+ $@"
	@go test -v -race -tags "$(BUILDTAGS) cgo" $(shell go list ${PROJECT}/... | grep -v vendor)

.PHONY: cover
cover:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

.PHONY: clean
clean:
	rm -f ${APP}

HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_LINT
	go get -u github.com/golang/lint/golint
endif
