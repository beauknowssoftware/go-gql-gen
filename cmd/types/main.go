package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
)

var sortFlag = flag.Bool("sort", false, "")

func main() {
	flag.Parse()

	schemaBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("failed to read schema from stdin: %v\n", err)
		os.Exit(1)
	}
	schema := string(schemaBytes)

	p := parse.New(parse.NewLexer(schema))
	ast, perr := p.Parse()
	if perr != nil {
		fmt.Printf("failed to parse schema: %v (%v)\n", perr.Error, perr.Token)
		os.Exit(1)
	}

	types := make([]string, 0)
	parse.Traverse(ast, func(n parse.Node) bool {
		if tdn, ok := n.(parse.TypeDefNode); ok {
			types = append(types, tdn.Name)
			return false
		}
		return true
	})

	if *sortFlag {
		sort.Slice(types, func(i, j int) bool {
			return types[i] < types[j]
		})
	}

	for _, n := range types {
		fmt.Println(n)
	}
}
