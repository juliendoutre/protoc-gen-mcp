package main

import (
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer(
        "{{ .Name }}",
        "{{ .Version }}",
    )

    if err := server.ServeStdio(s); err != nil {
        panic(err)
    }
}
