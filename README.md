# protoc-gen-mcp

## Prerequisites

```shell
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Run the example

```shell
go install .
protoc ./example/api.proto --mcp_out=./example/ --go_out=./example/ --go-grpc_out=./example/ -I./example/ -I./protos/
```

## Development

## Regenerate the extension protobuf

```shell
protoc ./protos/extension.proto --go_out=. --go-grpc_out=.
```

## References

- https://clement-jean.github.io/writing_protoc_plugins/
- https://github.com/theluckiestsoul/protoc-gen-structtag/blob/master/cmd/protoc-gen-structtag/main.go
- https://github.com/mark3labs/mcp-go/
