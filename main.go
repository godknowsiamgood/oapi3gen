package main

import (
	"api/oapi3"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func PathToMethodName(path string, method string) string {
	b := strings.Builder{}

	switch method {
	case "get":
		b.WriteString("Get")
	case "post":
		b.WriteString("Post")
	case "delete":
		b.WriteString("Delete")
	case "update":
		b.WriteString("Update")
	}

	r, _ := regexp.Compile("(\\w+)")

	path = strings.ReplaceAll(path, "_", " ")
	parts := r.FindAllString(path, 10)
	for _, p := range parts {
		b.WriteString(strings.Title(p))
	}

	return b.String()
}

func toColumnParametersPath(path string) string {
	path = strings.ReplaceAll(path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")
	return path
}

func templateMap(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func log(format string, vars ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", vars...)
}

func main() {
	if len(os.Args) == 1 {
		log("yml spec file should be provided")
		return
	}
	specFileName := os.Args[1]

	var genFileOutput string
	if len(os.Args) == 3 {
		genFileOutput = os.Args[2]
	}

	doc, err := openapi3.NewLoader().LoadFromFile(specFileName)
	if err != nil {
		log("schema parsing failed: %v", err)
		return
	}

	if err := doc.Validate(context.TODO()); err != nil {
		log("schema validation failed: %v", err)
		return
	}

	s := oapi3.Spec{}

	content, _ := ioutil.ReadFile(specFileName)
	if err := yaml.Unmarshal(content, &s); err != nil {
		log("schema parsing failed: %v", err)
		return
	}

	t, err := template.New("echo.tmpl").Funcs(template.FuncMap{
		"pathToMethodName":       PathToMethodName,
		"toCamel":                strcase.ToCamel,
		"toUpper":                strings.ToUpper,
		"toColumnParametersPath": toColumnParametersPath,
		"dict":                   templateMap,
	}).ParseFiles("echo.tmpl")
	if err != nil {
		log("failed to compile template: %v", err)
		return
	}

	sb := new(bytes.Buffer)

	if err := t.Execute(sb, s); err != nil {
		log("code generation failed: %v", err)
		return
	}

	formattedSource, err := format.Source(sb.Bytes())
	if err != nil {
		log("code formatting failed: %v", err)
		return
	}

	if genFileOutput != "" {
		genDirOutput := filepath.Dir(genFileOutput)
		if err := os.MkdirAll(genDirOutput, 0750); err != nil {
			log("saving generated code failed: %v", err)
		}
		if err := ioutil.WriteFile(genFileOutput, formattedSource, 0755); err != nil {
			log("saving generated code failed: %v", err)
			return
		}
	} else {
		if _, err := fmt.Fprintf(os.Stdout, "%s", formattedSource); err != nil {
			log("%v", err)
		}
	}
}
