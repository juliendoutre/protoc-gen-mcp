# protoc-gen-mcp

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
