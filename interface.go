package main

import (
	"github.com/godknowsiamgood/oapi3gen/spec"
	"strings"
	"text/template"
)

type Server interface {
	TemplateFunctions() template.FuncMap
	OperationParameterTags(param spec.Parameter) string
	FieldTags(context string, name string, field spec.Schema, parent spec.Schema) string
}

type DefaultServer struct{}

func (s DefaultServer) OperationParameterTags(_ spec.Parameter) string {
	return ""
}

func (s DefaultServer) TemplateFunctions() template.FuncMap {
	return nil
}

func (s DefaultServer) FieldTags(context string, name string, field spec.Schema, parent spec.Schema) string {
	tags := strings.Builder{}

	tags.WriteString("`json:\"" + name)

	if parent.IsFieldOptional(name) && context != spec.PropertiesContextParameters && context != spec.PropertiesContextRequestBody {
		tags.WriteString(",omitempty")
	}

	tags.WriteString("\"`")

	return tags.String()
}
