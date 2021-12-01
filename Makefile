PROJECT_NAME ?= api-example

build: ## Build main server
	@go build -ldflags="-X main.APIVersion=${API_VERSION}" -o bin/${PROJECT_NAME} -v

watch: ## Start compile deamon
	@CompileDaemon -build="make" -command="./bin/api-example" \
		-exclude-dir=.git \
		-exclude-dir=data