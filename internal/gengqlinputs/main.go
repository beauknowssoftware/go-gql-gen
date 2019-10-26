package gengqlinputs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/beauknowssoftware/go-gql-gen/internal/parse"
)

var packageFlag = flag.String("package", "", "")

func hasResolveDirective(fn parse.FieldNode) bool {
	for _, n := range fn.Directives {
		if dn, ok := n.(parse.DirectiveNode); ok && dn.Name == "resolve" {
			return true
		}
	}
	return false
}

func Run() {
	flag.Parse()

	pkg, hasPkgEnv := os.LookupEnv("GOPACKAGE")
	if !hasPkgEnv && *packageFlag == "" {
		fmt.Fprint(os.Stderr, "either GOPACKAGE environment variable or package flag must be set\n")
		os.Exit(1)
	} else if !hasPkgEnv {
		pkg = *packageFlag
	}

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

	fmt.Printf("package %v\n\n", pkg)
	fmt.Println("type ID string")
	fmt.Println()
	parse.Traverse(rnode, func(n parse.Node) bool {
		if tdn, ok := n.(parse.TypeDefNode); ok {
			if !tdn.Input {
				return false
			}
			fmt.Printf("type %v struct {\n", tdn.Name)
			parse.Traverse(tdn, func(n parse.Node) bool {
				if fn, ok := n.(parse.FieldNode); ok {
					tn := fn.Type.(parse.TypeNode)
					if tn.Name == "Query" {
						return false
					}
					if strings.HasSuffix(fn.Name, "Id") {
						prefix := strings.TrimSuffix(fn.Name, "Id")
						fmt.Printf("\t%vID", strings.Title(prefix))
					} else if fn.Name == "id" {
						fmt.Print("\tID")
					} else {
						fmt.Printf("\t%v", strings.Title(fn.Name))
					}
					switch tn.Name {
					case "String":
						fmt.Print(" string")
					default:
						if tn.Multiple {
							fmt.Printf(" []%v", tn.Name)
						} else if tn.Name == "ID" {
							fmt.Print(" ID")
						} else {
							fmt.Printf(" *%v", tn.Name)
						}
					}
					fmt.Printf(" `json:\"%v\"`", fn.Name)
					fmt.Println()
					return false
				}
				return true
			})
			fmt.Println("}")
			return false
		}
		return true
	})

	parse.Traverse(rnode, func(n parse.Node) bool {
		if tdn, ok := n.(parse.TypeDefNode); ok {
			if tdn.Input || tdn.Name == "Mutation" || tdn.Name == "Query" {
				return false
			}
			fmt.Println()
			fmt.Printf("type %v struct {\n", tdn.Name)
			parse.Traverse(tdn, func(n parse.Node) bool {
				if fn, ok := n.(parse.FieldNode); ok {
					tn := fn.Type.(parse.TypeNode)
					if tn.Name == "Query" {
						return false
					}
					if hasResolveDirective(fn) {
						fmt.Printf("\t%vLink\n", tn.Name)
						return false
					}
					if strings.HasSuffix(fn.Name, "Id") {
						prefix := strings.TrimSuffix(fn.Name, "Id")
						fmt.Printf("\t%vID", strings.Title(prefix))
					} else if fn.Name == "id" {
						fmt.Print("\tID")
					} else {
						fmt.Printf("\t%v", strings.Title(fn.Name))
					}
					switch tn.Name {
					case "String":
						if tn.Multiple {
							fmt.Print(" []string")
						} else {
							fmt.Print(" string")
						}
					case "Int":
						if tn.Multiple {
							fmt.Print(" []int")
						} else {
							fmt.Print(" int")
						}
					case "ID":
						fmt.Print(" ID")
					default:
						if tn.Multiple {
							fmt.Printf(" []%v", tn.Name)
						} else {
							fmt.Printf(" *%v", tn.Name)
						}
					}
					fmt.Printf(" `json:\"%v\"`", fn.Name)
					fmt.Println()
					return false
				}
				return true
			})
			fmt.Println("}")
			return false
		}
		return true
	})
}
