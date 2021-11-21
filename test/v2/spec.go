package v1

/**
    AUTOGENERATED. Please, do not edit.
**/

import (
	"net/http"
)

/* Components schemas */

type PetSchema struct {
	Id   int64   `json:"id"`
	Name string  `json:"name"`
	Tag  *string `json:"tag"`
}

/* Components responses */

/* Parameters */

type GetAaaParams struct {
}

type PostAaaParams struct {
	Test *int64
}

type PutAaaParams struct {
	Test *int64
}

type PostBbbParams struct {
}

type PostCccParams struct {
}

/* Requests bodies */

type PostAaaBody PetSchema

type PutAaaBody *PetSchema

type PostBbbBody []int64

/* Response objects */

/* Responses */

type Controller interface {
	GetAaa(params *GetAaaParams, req *http.Request, res http.ResponseWriter) int
	PostAaa(params *PostAaaParams, body *PostAaaBody, req *http.Request, res http.ResponseWriter) int
	PutAaa(params *PutAaaParams, body *PutAaaBody, req *http.Request, res http.ResponseWriter) int
	PostBbb(params *PostBbbParams, body *PostBbbBody, req *http.Request, res http.ResponseWriter) int
	PostCcc(params *PostCccParams, req *http.Request, res http.ResponseWriter) int
}
