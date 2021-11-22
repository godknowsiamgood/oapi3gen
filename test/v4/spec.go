package v1

/**
    AUTOGENERATED. Please, do not edit.
**/

import (
	"net/http"
)

/* Components schemas */

type InnerMapSchema map[string]interface{}

type InnerStructSchema struct {
	X string `json:"x,omitempty"`
	Y string `json:"y,omitempty"`
}

type OuterSchema struct {
	Inner1 InnerMapSchema     `json:"inner1"`
	Inner2 InnerMapSchema     `json:"inner2,omitempty"`
	Inner3 InnerStructSchema  `json:"inner3"`
	Inner4 *InnerStructSchema `json:"inner4,omitempty"`
	Z      string             `json:"z,omitempty"`
}

/* Components responses */

/* Parameters */

type GetTestParams struct {
	In1 InnerMapSchema
	In2 InnerMapSchema
	In3 *InnerStructSchema
	In4 InnerStructSchema
}

/* Requests bodies */

/* Response objects */

/* Responses */

type Controller interface {
	GetTest(params *GetTestParams, req *http.Request, res http.ResponseWriter) int
}