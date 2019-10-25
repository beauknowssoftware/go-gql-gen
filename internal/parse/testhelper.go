package parse

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGetDoc(t *testing.T, filename string) string {
	f, err := os.Open(path.Join("testdata", filename))
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		t.Fatalf("failed to open file %v\n%v", filename, err)
	}

	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read file data %v", err)
	}
	return string(fileData)
}

func TestParse(t *testing.T, schema string) DocumentNode {
	l := NewLexer(schema)
	p := New(l)
	d, err := p.Parse()
	if err != nil {
		t.Fatalf("failed to parse %v", err)
	}
	return *d
}
