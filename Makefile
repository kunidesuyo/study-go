IMAGE_TAG = web
ifeq ($(WEB_SERVER),echo)
	OSPI_CMD = oapi-codegen --config=./adapter/controller/echo/config.yaml ./api/openapi.yaml
else
	OSPI_CMD = oapi-codegen --config=./adapter/controller/gin/config.yaml ./api/openapi.yaml
endif

generate-code-from-openapi: ## Generate code from openapi
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	$(OSPI_CMD)

external-up: ## Up external containers
	pushd ./build/docker && docker-compose up -d mysql && popd

external-down: ## Down external containers
	pushd ./build/docker && docker-compose down && popd

mysql-cli: ## Connect to mysql cli
	pushd ./build/docker && docker-compose run mysql-cli && popd

run: ## Run app
	export APP_ENV=development
	go run ./cmd/server/main.go

docker-build: ## Build image
	docker build --tag $(IMAGE_TAG) -f ./build/docker/Dockerfile .

docker-run: ## Run docker
	docker run -p 8080:8080 -i -t $(IMAGE_TAG)

docker-compose-up: docker-build ## Run docker compose up
	pushd ./build/docker && docker-compose up -d --wait mysql web swagger-ui && popd

docker-compose-down: ## Run docker compose down
	pushd ./build/docker && docker-compose down && popd

unittest: ## Run unittest
	go clean -testcache
	go test -v ./adapter/... ./usecase/... ./entity/...

test-cover: ## Run test cover
	go test -coverprofile=coverage.out `go list ./... | grep -v /integration` && go tool cover -html=coverage.out

integration-test: generate-code-from-openapi docker-build ## Run integration test
	export APP_ENV=integration
	-pushd ./build/docker && docker-compose up -d --wait mysql web && popd
	-go clean -testcache
	-go test -v ./integration/...
	+pushd ./build/docker && docker-compose down && popd

lint: ## Run lint to show the diff
	# Install with `brew install golangci-lint` on Mac
	golangci-lint run

vet: ## Run go vet to show the diff
	go vet ./...

gofmt: ## Run gofmt to show the diff
	gofmt -d .

goimports: ## Run goimports to show the diff
	goimports -d .

prettier: vet gofmt goimports lint ## Show the code that needs to be modified

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $1, $2}'