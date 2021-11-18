package oapi3

import (
	"fmt"
	"github.com/Masterminds/semver"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Ref string

func (r Ref) IsSet() bool {
	return r != ""
}

func (r Ref) IsGenericError() bool {
	return r == "#/components/responses/Error"
}

func (r Ref) GetName() string {
	_, name := r.GetFullName()
	return name
}

func (r Ref) GetFullName() (string, string) {
	reg, _ := regexp.Compile("#/components/(schemas|responses)/(.+)")
	res := reg.FindAllStringSubmatch(string(r), 2)
	return res[0][1], res[0][2]
}

func (r Ref) GetTypeName() string {
	typ, name := r.GetFullName()

	switch typ {
	case "schemas":
		return name + "Schema"
	case "responses":
		return name + "Response"
	}

	return "Unknown"
}

type Type string

func (t Type) IsObject() bool {
	return t == "object"
}
func (t Type) IsArray() bool {
	return t == "array"
}
func (t Type) IsPrimitive() bool {
	return t == "string" || t == "number" || t == "boolean" || t == "integer"
}

func (t Type) GetTypeName() string {
	switch t {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "integer":
		return "int64"
	case "boolean":
		return "bool"
	}
	return "<unknown_type>"
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Parameter struct {
	In          string `yaml:"in"`
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Schema      Schema `yaml:"schema"`
}

func (p Parameter) IsRequired() bool {
	return p.Required || p.In == "path"
}

func (p Parameter) HasDefault() bool {
	return p.Schema.Default != nil
}

type Schema struct {
	AllOf []Schema `yaml:"allOf"`

	Ref                  Ref               `yaml:"$ref"`
	Description          string            `yaml:"description"`
	Type                 Type              `yaml:"type"`
	Minimum              *int              `yaml:"minimum"`
	Maximum              *int              `yaml:"maximum"`
	MinimumLength        *int              `yaml:"minLength"`
	MaximumLength        *int              `yaml:"maxLength"`
	Default              *string           `yaml:"default"`
	Required             []string          `yaml:"required"`
	Enum                 []string          `yaml:"enum"`
	Items                *Schema           `yaml:"items"`
	Properties           map[string]Schema `yaml:"properties"`
	AdditionalProperties interface{}       `yaml:"additionalProperties"`
}

func (s Schema) GetDefaultsTags() string {
	if s.Default == nil {
		return ""
	}

	return fmt.Sprintf("default:\"%v\"", *s.Default)
}

func (s Schema) GetValidationTags(isRequired bool) string {
	var tags []string

	if isRequired {
		tags = append(tags, "required")
	}
	if len(s.Enum) > 0 {
		tags = append(tags, "oneof="+strings.Join(s.Enum, " "))
	}
	if s.Minimum != nil {
		tags = append(tags, "min="+strconv.Itoa(*s.Minimum))
	}
	if s.Maximum != nil {
		tags = append(tags, "max="+strconv.Itoa(*s.Maximum))
	}
	if s.MinimumLength != nil {
		tags = append(tags, "min="+strconv.Itoa(*s.MinimumLength))
	}
	if s.MaximumLength != nil {
		tags = append(tags, "max="+strconv.Itoa(*s.MaximumLength))
	}
	return strings.Join(tags, ",")
}
func (s Schema) ExpandAllOf() {
	if len(s.AllOf) == 0 {
		return
	}

	s.Type = "object"

	if s.Properties == nil {
		s.Properties = make(map[string]Schema)
	}

	for _, e := range s.AllOf {
		for k, v := range e.Properties {
			s.Properties[k] = v
		}
	}
}

func (s Schema) GetSerializedEnums(asStrings bool) string {
	if len(s.Enum) > 0 {
		b := strings.Builder{}
		if asStrings {
			b.WriteString("[]string{")
		} else {
			b.WriteString("[]int64{")
		}

		for _, e := range s.Enum {
			if asStrings {
				b.WriteString("\"")
			}
			b.WriteString(e)
			if asStrings {
				b.WriteString("\"")
			}
			b.WriteString(",")
		}
		b.WriteString("}")
		return b.String()
	} else {
		return "nil"
	}
}

func (s Schema) IsFieldOptional(fieldName string) bool {
	for _, r := range s.Required {
		if fieldName == r {
			return false
		}
	}
	return true
}

func (s Schema) IsSet() bool {
	return s.Ref.IsSet() || len(s.AllOf) > 0 || s.Type != ""
}
func (s Schema) GetDefault() string {
	if s.Default == nil {
		return "nil"
	} else {
		if s.IsNumeric() {
			return "intPointer(" + (*s.Default) + ")"
		} else {
			return "stringPointer(\"" + (*s.Default) + "\")"
		}
	}
}
func (s Schema) GetMinimum() string {
	if !s.IsNumeric() || s.Minimum == nil {
		return "nil"
	} else {
		return "intPointer(" + strconv.Itoa(*s.Minimum) + ")"
	}
}
func (s Schema) GetMaximum() string {
	if !s.IsNumeric() || s.Maximum == nil {
		return "nil"
	} else {
		return "intPointer(" + strconv.Itoa(*s.Maximum) + ")"
	}
}
func (s Schema) GetMinimumLength() string {
	if s.IsNumeric() || s.MinimumLength == nil {
		return "nil"
	} else {
		return "intPointer(" + strconv.Itoa(*s.MinimumLength) + ")"
	}
}
func (s Schema) GetMaximumLength() string {
	if s.IsNumeric() || s.MaximumLength == nil {
		return "nil"
	} else {
		return "intPointer(" + strconv.Itoa(*s.MaximumLength) + ")"
	}
}
func (s Schema) IsNumeric() bool {
	return s.Type == "number" || s.Type == "integer"
}
func (s Schema) GetGoType() string {
	if s.Ref != "" {
		return s.Ref.GetTypeName()
	}

	b := strings.Builder{}

	if s.Type == "array" {
		b.WriteString("[]")
		if s.Items.Ref != "" {
			b.WriteString(s.Items.Ref.GetTypeName())
		} else {
			b.WriteString(s.Items.Type.GetTypeName())
		}
	} else {
		b.WriteString(s.Type.GetTypeName())
	}

	return b.String()
}

type Response struct {
	Ref         Ref    `yaml:"$ref"`
	Description string `yaml:"description"`
	Content     struct {
		JSON struct {
			Schema Schema `yaml:"schema"`
		} `yaml:"application/json"`
	} `yaml:"content"`
}

func (r Response) IsInlineStructType() bool {
	if r.Ref.IsSet() || r.Content.JSON.Schema.Ref.IsSet() {
		return false
	}

	if r.Content.JSON.Schema.Type.IsObject() && r.Content.JSON.Schema.AdditionalProperties == nil {
		return true
	}

	return false
}

func (r Response) IsEmpty() bool {
	return !r.Ref.IsSet() && !r.Content.JSON.Schema.Ref.IsSet() && r.Content.JSON.Schema.Type == ""
}

type OperationRequestBody struct {
	IsRequired bool `yaml:"required"`
	Content    struct {
		JSON struct {
			Schema Schema `yaml:"schema"`
		} `yaml:"application/json"`
	} `yaml:"content"`
}

type Operation struct {
	OperationId  string               `yaml:"operationId"`
	Parameters   []Parameter          `yaml:"parameters"`
	Summary      string               `yaml:"summary"`
	Description  string               `yaml:"description"`
	Responses    map[string]Response  `yaml:"responses"`
	RequestBody  OperationRequestBody `yaml:"requestBody"`
	XMiddlewares []string             `yaml:"x-middlewares"`
}

func (op Operation) HasRequestBody() bool {
	return op.RequestBody.Content.JSON.Schema.IsSet()
}

func (op Operation) IsAllEmptyResponses() bool {
	for code, r := range op.Responses {
		if code == "default" && r.Ref.IsGenericError() {
			continue
		}
		if !r.IsEmpty() {
			return false
		}
	}
	return true
}

func (op Operation) GetResponses() map[string]Response {
	responses := make(map[string]Response)
	for status, response := range op.Responses {
		if _, err := strconv.Atoi(status); err != nil {
			continue
		}
		if !response.Ref.IsSet() && !response.Content.JSON.Schema.IsSet() {
			continue
		}
		responses[status] = response
	}
	return responses
}

type Components struct {
	Schemas   map[string]Schema   `yaml:"schemas"`
	Responses map[string]Response `yaml:"responses"`
}

type PathOperation map[string]Operation

type Spec struct {
	Swagger    string                   `yaml:"swagger"`
	Info       Info                     `yaml:"info"`
	BasePath   string                   `yaml:"basePath"`
	Paths      map[string]PathOperation `yaml:"paths"`
	Components Components               `yaml:"components"`
}

func (s Spec) GetPackageName() string {
	v, err := semver.NewVersion(s.Info.Version)
	if err != nil {
		return "v1"
	} else {
		return "v" + strconv.Itoa(int(v.Major()))
	}
}

func (s Spec) HasGenericErrorResponse() bool {
	_, has := s.Components.Responses["Error"]
	return has
}

func (s Spec) traverseSchema(schema Schema, cb func(Schema) bool) bool {
	var traverse func(schema Schema) bool
	traverse = func(schema Schema) bool {
		if cb(schema) == true {
			return true
		}

		if schema.Ref.IsSet() {
			if refSchema, ok := s.Components.Schemas[schema.Ref.GetName()]; ok {
				return traverse(refSchema)
			}
		}

		return false
	}
	return traverse(schema)
}

func (s Spec) IsNillableSchema(schema Schema) bool {
	return s.traverseSchema(schema, func(schema Schema) bool {
		if schema.AdditionalProperties != nil {
			// will be generated map type
			return true
		}
		if schema.Type.IsArray() {
			// will be generated slice type
			return true
		}
		return false
	})
}

func (s Spec) IsOmmitableSchema(schema Schema) bool {
	return s.traverseSchema(schema, func(schema Schema) bool {
		if schema.AdditionalProperties != nil {
			// will be generated map type
			return true
		}
		if schema.Type.IsArray() {
			// will be generated slice type
			return true
		}
		if schema.Type.IsPrimitive() {
			// will be generated slice type
			return true
		}
		return false
	})
}

func (s Spec) GetUnderlyingSchema(ref Ref) Schema {
	component, name := ref.GetFullName()
	if component == "schemas" {
		schema := s.Components.Schemas[name]
		if schema.Ref.IsSet() {
			return s.GetUnderlyingSchema(schema.Ref)
		} else {
			return schema
		}
	} else if component == "responses" {
		response := s.Components.Responses[name]
		if response.Ref.IsSet() {
			return s.GetUnderlyingSchema(response.Ref)
		} else {
			return response.Content.JSON.Schema
		}
	}

	return Schema{}
}

func (s Spec) GetAllMiddlewareNames() []string {
	middlewaresMap := make(map[string]bool)
	for _, operations := range s.Paths {
		for _, operation := range operations {
			for _, middlewareName := range operation.XMiddlewares {
				middlewaresMap[middlewareName] = true
			}
		}
	}

	var middlewares []string
	for name := range middlewaresMap {
		middlewares = append(middlewares, name)
	}

	sort.Strings(middlewares)

	return middlewares
}
