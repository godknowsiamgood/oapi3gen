package oapi3

import (
	"github.com/Masterminds/semver"
	"regexp"
	"strconv"
	"strings"
)

type Ref string

func (r Ref) IsSet() bool {
	return r != ""
}

func (r Ref) GetName() (string, string) {
	reg, _ := regexp.Compile("#/components/(schemas|responses)/(.+)")
	res := reg.FindAllStringSubmatch(string(r), 2)
	return res[0][1], res[0][2]
}

func (r Ref) GetTypeName() string {
	typ, name := r.GetName()

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

func (t Type) GetTypeName() string {
	switch t {
	case "string":
		return "string"
	case "number":
		return "float32"
	case "integer":
		return "int32"
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
			b.WriteString("[]int32{")
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

func (s Schema) IsFieldRequired(fieldName string) bool {
	for _, r := range s.Required {
		if fieldName == r {
			return true
		}
	}
	return false
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

type Operation struct {
	Parameters  []Parameter         `yaml:"parameters"`
	Summary     string              `yaml:"summary"`
	Description string              `yaml:"description"`
	Responses   map[string]Response `yaml:"responses"`
}

func (op Operation) IsAllEmptyResponses() bool {
	for _, r := range op.Responses {
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

func (s Spec) HasComponentErrorResponse() bool {
	_, has := s.Components.Responses["Error"]
	return has
}

func (s Spec) IsNillableSchema(schema Schema) bool {
	var check func(schema Schema) bool
	check = func(schema Schema) bool {
		if schema.AdditionalProperties != nil {
			// will be generated map type
			return false
		}

		if schema.Type.IsArray() {
			// will be generated slice type
			return false
		}

		if schema.Ref.IsSet() {
			_, name := schema.Ref.GetName()
			if refSchema, ok := s.Components.Schemas[name]; ok {
				return check(refSchema)
			}
		}

		return true
	}
	return check(schema)
}
