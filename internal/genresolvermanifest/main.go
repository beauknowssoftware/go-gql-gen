package genresolvermanifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beauknowssoftware/go-gql-gen/internal/parse"
)

func hasResolveDirective(fn parse.FieldNode) bool {
	for _, n := range fn.Directives {
		if dn, ok := n.(parse.DirectiveNode); ok && dn.Name == "resolve" {
			return true
		}
	}
	return false
}

func Run() {
	schemaBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schema from stdin: %v\n", err)
		os.Exit(1)
	}
	schema := string(schemaBytes)

	p := parse.New(parse.NewLexer(schema))
	rnode, perr := p.Parse()
	if perr != nil {
		fmt.Fprintf(os.Stderr, "failed to parse schema: %v (%v)\n", perr.Error, perr.Token)
		os.Exit(1)
	}

	type Entry struct {
		Type  string `json:"type"`
		Field string `json:"field"`
	}
	entries := make([]Entry, 0)
	parse.Traverse(rnode, func(n parse.Node) bool {
		if tdn, ok := n.(parse.TypeDefNode); ok {
			for _, n := range tdn.Fields {
				fn := n.(parse.FieldNode)
				if hasResolveDirective(fn) {
					entries = append(entries, Entry{
						Type:  tdn.Name,
						Field: fn.Name,
					})
				}
			}
			return false
		}
		return true
	})

	d, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal manifest: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(d)
}
