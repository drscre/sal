package main

import (
	"io/ioutil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update bool = false

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode("repo", "github.com/go-gad/sal/examples/bookstore", []string{"Store"})
	if err != nil {
		t.Fatalf("Failed to generate a code: %+v", err)
	}

	//t.Logf("\n%s", string(code))
	if update {
		if err = ioutil.WriteFile("../examples/bookstore/repo/sal_client.go", code, 0666); err != nil {
			t.Fatalf("failed to write file: %+v", err)
		}
	}
	expCode, err := ioutil.ReadFile("../examples/bookstore/repo/sal_client.go")
	if string(expCode) != string(code) {
		t.Error("generated code is not equal to expected")
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expCode), string(code), true)
		t.Log(dmp.DiffPrettyText(diffs))
	}
}
