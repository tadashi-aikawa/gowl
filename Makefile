build: ## Static build
	@go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"'