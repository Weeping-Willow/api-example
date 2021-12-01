PROJECT_NAME ?= api-example

build: ## Build main server
	@go build -o bin/${PROJECT_NAME}

watch: ## Start compile deamon
	@CompileDaemon -build="make" -command="./bin/api-example" \
		-exclude-dir=.git \
		-exclude-dir=data