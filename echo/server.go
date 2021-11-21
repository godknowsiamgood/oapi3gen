package echo

import (
	"fmt"
	"github.com/godknowsiamgood/oapi3gen/spec"
	"strconv"
	"strings"
	"text/template"
)

type Server struct {
	Spec spec.Spec
}

func toColumnParametersPath(path string) string {
	path = strings.ReplaceAll(path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")
	return path
}

func getDefaultStatusCode(pattern string) string {
	if pattern == "default" {
		return "200"
	} else if strings.Contains(pattern, "X") {
		return strings.ReplaceAll(pattern, "X", "0")
	} else {
		return pattern
	}
}

func (s *Server) TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"toColumnParametersPath": toColumnParametersPath,
		"toUpper":                strings.ToUpper,
		"getDefaultStatusCode":   getDefaultStatusCode,
	}
}

func (s *Server) OperationParameterTags(param spec.Parameter) string {
	var tags []string

	// location
	switch param.In {
	case "path":
		tags = append(tags, "param:\""+param.Name+"\"")
	case "query":
		tags = append(tags, "query:\""+param.Name+"\"")
	}

	// defaults
	if param.Schema.Default != nil {
		tags = append(tags, fmt.Sprintf("default:\"%v\"", *param.Schema.Default))
	}

	// validation
	var validations []string
	if param.Required {
		validations = append(validations, "required")
	}
	schema := param.Schema
	if len(schema.Enum) > 0 {
		validations = append(validations, "oneof="+strings.Join(schema.Enum, " "))
	}
	if schema.Minimum != nil {
		validations = append(validations, "min="+strconv.Itoa(*schema.Minimum))
	}
	if schema.Maximum != nil {
		validations = append(validations, "max="+strconv.Itoa(*schema.Maximum))
	}
	if schema.MinimumLength != nil {
		validations = append(validations, "min="+strconv.Itoa(*schema.MinimumLength))
	}
	if schema.MaximumLength != nil {
		validations = append(validations, "max="+strconv.Itoa(*schema.MaximumLength))
	}
	if len(validations) == 0 {
		tags = append(tags, "validate:\""+strings.Join(validations, ",")+"\"")
	}

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
}
