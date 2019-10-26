package main_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func OpenFile(t *testing.T, filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open file %v", filename)
	}
	return f
}

func ReadFile(t *testing.T, filename string) string {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file %v", filename)
	}
	return string(d)
}

func Test_Main(t *testing.T) {
	tests := []string{
		"types",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			gqls := OpenFile(t, "testdata/"+name+".graphqls")

			cmd := exec.Command("gen-gql-types", "-package", "test")
			var outBuff, errBuff bytes.Buffer
			cmd.Stdin = gqls
			cmd.Stdout = &outBuff
			cmd.Stderr = &errBuff

			if err := cmd.Run(); err != nil {
				t.Fatalf("failed to run %v\n%v\n%v", err, outBuff.String(), errBuff.String())
			}

			expected := ReadFile(t, "testdata/"+name+".go.test")
			if diff := cmp.Diff(expected, outBuff.String()); diff != "" {
				t.Fatalf("mismatch (-expected,+got) %v", diff)
			}
		})
	}
}
