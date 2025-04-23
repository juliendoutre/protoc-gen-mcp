package main

import (
    "github.com/mark3labs/mcp-go/server"
    {{ if .Tools }}"github.com/mark3labs/mcp-go/mcp"{{ end }}
)

func main() {
    s := server.NewMCPServer(
        "{{ .Name }}",
        "{{ .Version }}",
    )

    {{ range .Tools }}
    {{ .Name }}Tool := mcp.NewTool("{{ .Name }}",
        {{ if .Description }}mcp.Description("{{ .Description }}"),{{ end }}
    )
    s.AddTool({{ .Name }}Tool, nil)
    {{ end }}

    if err := server.ServeStdio(s); err != nil {
        panic(err)
    }
}
