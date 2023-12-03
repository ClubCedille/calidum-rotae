default: setup

help:  	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

setup: 	## Setup command
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# gRPC protoc template
define grpc_template
	protoc -I./api --go_out=$(1) --go_opt=paths=source_relative \
    --go-grpc_out=$(1) --go-grpc_opt=paths=source_relative \
    $(2)
endef

grpc: 	## Generate gRPC files
	$(call grpc_template,./pkg/proto-gen/provider,api/provider.proto)
	$(call grpc_template,./pkg/proto-gen/email-provider,api/email_provider.proto)
	$(call grpc_template,./pkg/proto-gen/discord-provider,api/discord_provider.proto)
	
##########################################
####### Docker related commands ##########
##########################################
docker-compose:
	@docker-compose build --no-cache && docker-compose up

docker-compose-dev:
	@docker-compose -f docker-compose.dev.yml up --build
