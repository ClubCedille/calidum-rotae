default: setup

help:  	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

setup: 	## Setup command
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

grpc: ## Generate gRPC files
	for file in api/*.proto; do \
		name=$$(basename "$$file" .proto); \
		name=$${name//_/-}; \
		mkdir -p "./pkg/proto-gen/$$name"; \
		protoc -I./api \
			--go_out=./pkg/proto-gen/$$name --go_opt=paths=source_relative \
			--go-grpc_out=./pkg/proto-gen/$$name --go-grpc_opt=paths=source_relative \
			"$$file"; \
	done
##########################################
####### Docker related commands ##########
##########################################
docker-compose:
	@docker-compose build --no-cache && docker-compose up

docker-compose-dev:
	@docker-compose -f docker-compose.dev.yml up --build
