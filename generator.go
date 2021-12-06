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
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	r    = regexp.MustCompile(`\n(\n+)`)
	repl = []byte("\n")
)

func preprocess(b string) string {
	return string(r.ReplaceAll([]byte(b), repl))
}

func generate(yamlContent []byte, serverName string) ([]byte, error) {
	if isVerbose {
		log("Validating spec for yaml file (%v bytes)...", len(yamlContent))
	}

	doc, err := openapi3.NewLoader().LoadFromData(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("schema parsing failed: %v", err)
	}

	if err := doc.Validate(context.TODO()); err != nil {
		return nil, fmt.Errorf("schema validation failed: %v", err)
	}

	if isVerbose {
		log("Parsing spec...")
	}

	s := spec.Spec{}

	if err := yaml.Unmarshal(yamlContent, &s); err != nil {
		return nil, fmt.Errorf("schema parsing failed: %v", err)
	}

	if isVerbose {
		log("Parsed spec %s", s.Info.Title)
	}

	var server Server
	switch serverName {
	case "echo":
		server = &echo.Server{Spec: s}
	default:
		server = DefaultServer{}
	}

	if isVerbose {
		log("Loading templates...")
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

	baseTemplateContent = strings.Replace(baseTemplateContent, "{{/*server*/}}", serverTmplReplace, 1) //
	baseTemplateContent = strings.Replace(baseTemplateContent, "{{/*boilerplate*/}}", serverBoilerplateReplace, 1)

	var addedImports []string

	objectsContext := "default"

	templateFunctions := template.FuncMap{
		"operationId":             spec.OperationId,
		"toCamel":                 strcase.ToCamel,
		"toLowerCamel":            strcase.ToLowerCamel,
		"dict":                    templateMap,
		"isNillableSchema":        s.IsNillableSchema,
		"isOmmitableSchema":       s.IsOmmitableSchema,
		"isStruct":                s.IsStruct,
		"getUnderlyingSchema":     s.GetUnderlyingSchema,
		"hasGenericErrorResponse": s.HasGenericErrorResponse,
		"server":                  func() Server { return server },
		"getContext":              func() string { return objectsContext },
		"setContext": func(c string) string {
			objectsContext = c
			return ""
		},
		"addImport": func(imp string) string {
			addedImports = append(addedImports, imp)
			return ""
		},
	}

	serverTemplateFunctions := server.TemplateFunctions()
	for k, v := range serverTemplateFunctions {
		templateFunctions[k] = v
	}

	if isVerbose {
		log("Parsing templates...")
	}

	t, err := template.New("base.tmpl").Funcs(templateFunctions).Parse(baseTemplateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to compile template: %v", err)
	}

	sb := new(bytes.Buffer)

	if isVerbose {
		log("Generating code...")
	}

	if err := t.Execute(sb, s); err != nil {
		return nil, fmt.Errorf("code generation failed: %v", err)
	}

	sourceRaw := sb.Bytes()
	sourceRaw = renderImports(sourceRaw, addedImports)

	var formattedSource []byte
	if os.Getenv("debug") == "" {
		formattedSource, err = imports.Process("", sourceRaw, nil)
		if err != nil {
			return nil, fmt.Errorf("code formatting failed: %v", err)
		}
	} else {
		formattedSource = sourceRaw
	}

	return formattedSource, nil
}

func renderImports(sourceRaw []byte, imports []string) []byte {
	importsStr := strings.Builder{}
	for _, i := range imports {
		importsStr.WriteString("import \"" + i + "\"\n")
	}
	return []byte(strings.Replace(string(sourceRaw), "/*imports*/", importsStr.String(), 1))
}
