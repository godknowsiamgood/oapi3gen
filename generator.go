package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-yaml"
	"github.com/godknowsiamgood/oapi3gen/echo"
	"github.com/godknowsiamgood/oapi3gen/spec"
	"github.com/iancoleman/strcase"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func generate(yamlContent []byte, serverName string) ([]byte, error) {
	doc, err := openapi3.NewLoader().LoadFromData(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("schema parsing failed: %v", err)
	}

	if err := doc.Validate(context.TODO()); err != nil {
		return nil, fmt.Errorf("schema validation failed: %v", err)
	}

	s := spec.Spec{}

	if err := yaml.Unmarshal(yamlContent, &s); err != nil {
		return nil, fmt.Errorf("schema parsing failed: %v", err)
	}

	var server Server
	switch serverName {
	case "echo":
		server = &echo.Server{Spec: s}
	default:
		server = DefaultServer{}
	}

	baseTemplateFile, _ := ioutil.ReadFile("./base.tmpl")
	baseTemplateContent := string(baseTemplateFile)

	var serverTmplReplace string
	var serverBoilerplateReplace string
	if serverName != "" {
		serverTmpl, _ := ioutil.ReadFile("./" + serverName + "/server.tmpl")
		if len(serverTmpl) != 0 {
			serverTmplReplace = string(serverTmpl)
			serverBoilerplateReplace = "{{ template \"serverBoilerplate\" . }}"
		}
	}

	baseTemplateContent = strings.Replace(baseTemplateContent, "%server%", serverTmplReplace, 1)
	baseTemplateContent = strings.Replace(baseTemplateContent, "%boilerplate%", serverBoilerplateReplace, 1)

	isOmittingFields := false

	templateFunctions := template.FuncMap{
		"operationId":             spec.OperationId,
		"toCamel":                 strcase.ToCamel,
		"toLowerCamel":            strcase.ToLowerCamel,
		"dict":                    templateMap,
		"isNillableSchema":        s.IsNillableSchema,
		"isOmmitableSchema":       s.IsOmmitableSchema,
		"getUnderlyingSchema":     s.GetUnderlyingSchema,
		"hasGenericErrorResponse": s.HasGenericErrorResponse,
		"server":                  func() Server { return server },
		"isOmittingFields":        func() bool { return isOmittingFields },
		"setOmittingFields": func(set bool) string {
			isOmittingFields = set
			return ""
		},
	}

	serverTemplateFunctions := server.TemplateFunctions()
	for k, v := range serverTemplateFunctions {
		templateFunctions[k] = v
	}

	t, err := template.New("echo.tmpl").Funcs(templateFunctions).Parse(baseTemplateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to compile template: %v", err)
	}

	sb := new(bytes.Buffer)

	if err := t.Execute(sb, s); err != nil {
		return nil, fmt.Errorf("code generation failed: %v", err)
	}

	var formattedSource []byte
	if os.Getenv("debug") == "" {
		formattedSource, err = format.Source(sb.Bytes())
		if err != nil {
			return nil, fmt.Errorf("code formatting failed: %v", err)
		}
	} else {
		formattedSource = sb.Bytes()
	}

	return formattedSource, nil
}
