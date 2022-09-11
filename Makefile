.DEFAULT_GOAL := help

# protoc compiler command for generating grpc go code
PROTOC_GO=protoc --proto_path=./grpc/idl --go_out=./grpc/go --go-grpc_out=./grpc/go --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative

# capnp compiler command for generating capnp go code
#  option -I specifies import location for non-relative imports
#  option -o<lang>:<dir> specifies the language (go) and the directory of the go generated code
#  option --src-prefix prevents generating generating the code in the subdirectory capnp/idl 
CAPNP_GO=capnp compile "-I$(GOPATH)/src/capnproto.org/go/capnp/std" -ogo:./capnp/go/ --src-prefix=capnp/idl
.FORCE:


setup: .FORCE ## Go Get grpc and capnp modules needed for building
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


echoboot: .FORCE ## run the echoboot app
	go run pkg/echoboot/main.go

capnp: .FORCE  ## generate the capnp based service
	$(CAPNP_GO)  ./capnp/idl/echo.capnp
	go mod tidy


grpc: .FORCE  ## generate the grpc based service
	$(PROTOC_GO) grpc/idl/echo.proto
	go mod tidy


help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


run: ## invoke the service via grpc
	go run pkg/main.go