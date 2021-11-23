package main

import (
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {
	yamlContent, _ := ioutil.ReadFile("./test/v1/spec.yaml")
	goContent, _ := ioutil.ReadFile("./test/v1/spec.gocode")
	out, err := generate(yamlContent, "")
	if err != nil {
		t.Error(err)
	}
	if string(out) != string(goContent) {
		t.Errorf("failed")
	}

	yamlContent, _ = ioutil.ReadFile("./test/v1/spec.yaml")
	goContent, _ = ioutil.ReadFile("./test/v1/spec_echo.gocode")
	out, err = generate(yamlContent, "echo")
	if err != nil {
		t.Error(err)
	}
	if string(out) != string(goContent) {
		t.Errorf("failed")
	}
}
