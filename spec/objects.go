package spec

import (
	"github.com/Masterminds/semver"
	"github.com/iancoleman/strcase"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const PropertiesContextComponents string = "components"
const PropertiesContextParameters string = "parameters"
const PropertiesContextRequestBody string = "requestBody"

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

type Schema struct {
	AllOf []Schema `yaml:"allOf"`
	AnyOf []Schema `yaml:"anyOf"`

	Ref                  Ref               `yaml:"$ref"`
	Description          string            `yaml:"description"`
	Type                 Type              `yaml:"type"`
	Format               string            `yaml:"format"`
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

func (s Schema) HasDefault() bool {
	return s.Default != nil
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
	return s.Ref.IsSet() || len(s.AllOf) > 0 || len(s.AnyOf) > 0 || s.Type != ""
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
		switch s.Format {
		case "int32":
			b.WriteString("int32")
		case "int64":
			b.WriteString("int64")
		case "float":
			b.WriteString("float32")
		case "double":
			b.WriteString("float64")
		default:
			b.WriteString(s.Type.GetTypeName())
		}
	}

	return b.String()
}

type Response struct {
	Ref         Ref     `yaml:"$ref"`
	Description string  `yaml:"description"`
	Content     Content `yaml:"content"`
}

func (r Response) IsInlineStructType() bool {
	contentJSONSchema := r.Content.GetBindableParametersSchema()

	if r.Ref.IsSet() || contentJSONSchema.Ref.IsSet() {
		return false
	}

	if contentJSONSchema.Type.IsObject() && contentJSONSchema.AdditionalProperties == nil {
		return true
	}

	return false
}

func (r Response) IsEmpty() bool {
	return !r.Ref.IsSet() && !r.Content.GetBindableParametersSchema().Ref.IsSet() && r.Content.GetBindableParametersSchema().Type == ""
}

type Content map[string]struct {
	Schema Schema `yaml:"schema"`
}

func (c Content) GetBindableParametersSchema() Schema {
	for t, s := range c {
		if t == "application/json" || t == "multipart/form-data" {
			return s.Schema
		}
	}
	return Schema{}
}

func (c Content) IsParametrizedContent() bool {
	for t := range c {
		if t == "multipart/form-data" || t == "application/json" {
			return true
		}
	}
	return false
}

type OperationRequestBody struct {
	IsRequired bool    `yaml:"required"`
	Content    Content `yaml:"content"`
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

func (op Operation) HasRequestBodyBindableParameters() bool {
	if op.RequestBody.Content == nil {
		return false
	}
	return op.RequestBody.Content.GetBindableParametersSchema().IsSet()
}

func (op Operation) HasRequestBody() bool {
	return op.RequestBody.Content != nil
}

func (op Operation) HasRequestParametrizedBody() bool {
	if op.RequestBody.Content == nil {
		return false
	}
	return op.RequestBody.Content.IsParametrizedContent()
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
		if !response.Ref.IsSet() && !response.Content.GetBindableParametersSchema().IsSet() {
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

func (s Spec) IsStruct(schema Schema) bool {
	return s.traverseSchema(schema, func(schema Schema) bool {
		if schema.Type.IsObject() && schema.AdditionalProperties == nil {
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
			return response.Content.GetBindableParametersSchema()
		}
	}

	return Schema{}
}

func (s Spec) HasGenericErrorResponse() bool {
	_, has := s.Components.Responses["Error"]
	return has
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

func OperationId(path string, method string, operation Operation) string {
	if operation.OperationId != "" {
		return strcase.ToCamel(operation.OperationId)
	}

	b := strings.Builder{}

	switch method {
	case "get":
		b.WriteString("Get")
	case "post":
		b.WriteString("Post")
	case "put":
		b.WriteString("Put")
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
