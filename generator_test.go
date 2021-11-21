package main

import (
	"io/ioutil"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	for i := 1; i <= 2; i++ {
		yamlContent, _ := ioutil.ReadFile("./test/v" + strconv.Itoa(i) + "/spec.yaml")
		goContent, _ := ioutil.ReadFile("./test/v" + strconv.Itoa(i) + "/spec.go")
		out, err := generate(yamlContent, "")
		if err != nil {
			t.Error(err)
		}
		if string(out) != string(goContent) {
			t.Error(":(")
		}
	}
}
