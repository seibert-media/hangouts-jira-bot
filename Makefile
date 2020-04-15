-include .env

.PHONY: build

updateDebugger:
	wget -O files/go-cloud-debug https://storage.googleapis.com/cloud-debugger/compute-go/go-cloud-debug
	chmod 0755 files/go-cloud-debug


teamvault-config-dir-generator:
	@GO111MODULE=off go get github.com/bborbe/teamvault-utils/cmd/teamvault-config-dir-generator

generate: teamvault-config-dir-generator
	@test -f ${TEAMVAULT_SM} || echo "\nMissing file ${TEAMVAULT_SM} :\n{\n \"url\": \"https://teamvault.apps.seibert-media.net\",\n \"user\": \"mmustermann\",\n \"pass\": \"PASSWORT\"\\n}"
	@test -f ${TEAMVAULT_SM} || exit 1

	@rm -rf kubernetes-manifests/results
	teamvault-config-dir-generator \
		-teamvault-config="${TEAMVAULT_SM}" \
		-source-dir=kubernetes-manifests/templates \
		-target-dir=kubernetes-manifests/results \
		-logtostderr \
		-v=2

deploySecret: generate
	kubectl apply -n hangouts-jira-bot -f kubernetes-manifests/results
	@rm -rf kubernetes-manifests/results

# test entire repo
gotest:
	@go test -cover -race $(shell go list ./... | grep -v /vendor/)

test:
	@GO111MODULE=off go get github.com/onsi/ginkgo/ginkgo
	@ginkgo -r -race

# format entire repo (excluding vendor)
format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +
