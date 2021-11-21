package main

import (
	"github.com/godknowsiamgood/oapi3gen/spec"
	"text/template"
)

type Server interface {
	TemplateFunctions() template.FuncMap
	OperationParameterTags(param spec.Parameter) string
}

type DefaultServer struct{}

func (s DefaultServer) OperationParameterTags(_ spec.Parameter) string {
	return ""
}

func (s DefaultServer) TemplateFunctions() template.FuncMap {
	return nil
}
