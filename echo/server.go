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

func (s *Server) FieldTags(context string, name string, field spec.Schema, parent spec.Schema) string {
	var tags []string

	if context == spec.PropertiesContextComponents {
		omitempty := ""
		if parent.IsFieldOptional(name) {
			omitempty = ",omitempty"
		}
		tags = append(tags, "json:\""+name+omitempty+"\"")
	} else if context == spec.PropertiesContextRequestBody {
		tags = append(tags, "form:\""+name+"\"")
		tags = append(tags, getValidationAndDefaultTagsForInputSchema(field, !parent.IsFieldOptional(name))...)
	}

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
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

	tags = append(tags, getValidationAndDefaultTagsForInputSchema(param.Schema, param.Required)...)

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
}

func getValidationAndDefaultTagsForInputSchema(schema spec.Schema, isRequired bool) []string {
	var tags []string

	// defaults
	if schema.Default != nil {
		tags = append(tags, fmt.Sprintf("default:\"%v\"", *schema.Default))
	}

	// validation
	var validations []string
	if isRequired {
		validations = append(validations, "required")
	}

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
	if len(validations) > 0 {
		tags = append(tags, "validate:\""+strings.Join(validations, ",")+"\"")
	}

	return tags
}
