###################### //S/M Makefile ######################
#
# Parts of this makefile are based upon github.com/kolide/kit
#

NAME         := hangouts-jira-bot
REPO         := seibert-media
GIT_HOST     := github.com
REGISTRY     := quay.io
IMAGE        := seibertmedia/hangouts-jira-bot

PATH := $(GOPATH)/bin:$(PATH)
TOOLS_DIR := cmd
VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)

NAMESPACE	?= $(NAME)
DOCKER_CACHE ?= --no-cache
DOCKER_TAGS := -t $(REGISTRY)/$(IMAGE):$(VERSION) -t $(REGISTRY)/$(IMAGE):latest

-include .env

include helpers/make_version

.PHONY: build

### MAIN STEPS ###

all: test install run

# install required tools and dependencies
deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint
	go get -u github.com/haya14busa/goverage
	go get -u github.com/kisielk/errcheck
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/schrej/godacov
	go get -u golang.org/x/tools/cmd/goimports

updateDebugger:
	wget -O files/go-cloud-debug https://storage.googleapis.com/cloud-debugger/compute-go/go-cloud-debug
	chmod 0755 files/go-cloud-debug

# test entire repo
gotest:
	@go test -cover -race $(shell go list ./... | grep -v /vendor/)

test:
	@go get github.com/onsi/ginkgo/ginkgo
	@ginkgo -r -race

# install passed in tool project
install:
	GOBIN=$(GOPATH)/bin go install cmd/$(NAME)/*.go

# run tool
run:
	$(NAME)

# format entire repo (excluding vendor)
format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

# build binary for docker image
buildgo: .pre-build
	CGO_ENABLED=0 GOOS=linux go build -a -o app -ldflags ${KIT_VERSION} -gcflags='-N -l' ./cmd/$(NAME)

# build ui
buildui: .pre-build
	npm run --prefix ui build --env=prod

# build docker image
build:
	gcloud debug source gen-repo-info-file --output-directory files
	@docker build $(DOCKER_CACHE) --rm=true $(DOCKER_TAGS) \
	--build-arg VERSION=$(VERSION) \
	--build-arg GIT_HOST=$(GIT_HOST) \
	--build-arg REPO=$(REPO) \
	--build-arg NAME=$(NAME) \
	-f Dockerfile .

# run specified tool from code
dev: generate
	@go run -ldflags $(KIT_VERSION) cmd/$(NAME)/*.go \
	-debug

# build the docker image
docker: build

# upload the docker image
upload:
	docker push $(REGISTRY)/$(IMAGE)

# clean build results and delete all images
clean:
	rm -rf build
	docker rmi -f $(shell docker images -q --filter=reference=$(REGISTRY)/$(IMAGE)*)

version:
	@echo $(VERSION)

# create build dir
.pre-build:
	@mkdir -p build

# helper to build new image and kick existing pod
update-deployment: docker upload clean restart-deployment

# delete existing pod to force imagePull (if latest)
restart-deployment:
	kubectl delete po -n $(NAMESPACE) -lapp=$(NAME)

# go quality checks
check: format lint vet

# vet entire repo (excluding vendor)
vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

# lint entire repo (excluding vendor)
lint:
	@go get github.com/golang/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

# errcheck entire repo (excluding vendor)
errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write)' $(shell go list ./... | grep -v /vendor/)

cover:
	@go get github.com/haya14busa/goverage
	@go get github.com/schrej/godacov
	goverage -v -coverprofile=coverage.out $(shell go list ./... | grep -v /vendor/)

generate:
	@go get github.com/maxbrunsfeld/counterfeiter
	@go generate ./...
