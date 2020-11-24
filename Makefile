APP_VERSION ?=v0.5.2
APP_ID      ?=git-ops
APP_PORT    ?=8080
IMAGE_OWNER ?=$(shell git config --get user.username)

.PHONY: all
all: help

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

.PHONY: run
run: tidy ## Runs uncompiled code
	go run main.go

.PHONY: dapr
dapr: tidy ## Runs uncompiled code in Dapr
	dapr run \
	 --app-id gitops \
	 --app-port $(APP_PORT) \
	 --app-protocol http \
	 --components-path ./component \
	 --log-level debug \
	 go run main.go

.PHONY: image
image: tidy ## Builds and publish image 
	docker build \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%T-UTC") \
		-t ghcr.io/$(IMAGE_OWNER)/$(APP_ID):$(APP_VERSION) \
		.
	docker push ghcr.io/$(IMAGE_OWNER)/$(APP_ID):$(APP_VERSION)

.PHONY: lint
lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

.PHONY: sync
sync: ## Adds, commits, and pushes the local changes 
	git add .
	git commit -m "new greetings"
	git push --all

tag: ## Creates release tag 
	git tag $(APP_VERSION)
	git push origin $(APP_VERSION)

.PHONY: clean
clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
