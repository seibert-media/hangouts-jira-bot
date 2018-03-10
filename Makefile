###################### //S/M Makefile ######################
#
# This Makefile is used to manage the command-line template
# All possible tools have to reside under their respective folders in cmd/
# and are being autodetected.
# 'make full' would then process them all while 'make toolname' would only
# handle the specified one(s).
# Edit this file with care, as it is also being used by our CI/CD Pipeline
# For usage information check README.md
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

-include .env

include helpers/make_version
include helpers/make_gohelpers
include helpers/make_dockerbuild
include helpers/make_db

### MAIN STEPS ###

default:

# install required tools and dependencies
deps:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/dep/cmd/dep
ifdef BUILD_DEB
	go get -u github.com/bborbe/debian_utils/bin/create_debian_package
endif
	dep ensure

# install passed in tool project
install:
	GOBIN=$(GOPATH)/bin go install $(TOOLS_DIR)/*.go

# build passed in tool project
build: .pre-build
	@GOBIN=$(GOPATH)/bin go build -i -o build/$(NAME) -ldflags ${KIT_VERSION} ./$(TOOLS_DIR)/

# run specified tool binary
run: build
	@./build/$(NAME)

# run specified tool from code
dev:
	@go run -ldflags ${KIT_VERSION} $(TOOLS_DIR)/*.go -version=false -debug
# build the docker image
docker: build-in-docker build-image

# upload the docker image
upload:
	docker push $(REGISTRY)/$(IMAGE)

### HELPER STEPS ###

# clean local vendor folder
clean:
	rm -rf build
	docker rmi -f $(shell docker images -q --filter=reference=$(REGISTRY)/$(IMAGE)*)

build-docker-bin:
	GOBIN=$(GOPATH)/bin CGO_ENABLED=0 GOOS=linux go build -i -o build/app -ldflags ${KIT_VERSION_DOCKER} -a -installsuffix cgo ./$(TOOLS_DIR)/

.pre-build:
	@mkdir -p build

.build-all:
	make full build

update-deployment: docker upload clean restart-deployment

restart-deployment:
	kubectl delete po -n bot -lapp=$(NAME)
