# oapi3gen
Simple OpenAPI 3 code generator with zero or minimum server boilerplate

# Features
* Strongly typed parameters and responses
* No server boilerplate if not needed
* Minimum server boilerplate - just routes and parameter validation
* Easy middlewares with custom `x-middlewares` object
* Easy error responses - just add `components/schemas/Error` object

# Code generation
Program generates only one file with strongly types for request parameters, request bodies and responses.
Also it generates interface type for controller that you may implement.

E.g. specification like
```
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      tags:
        - pets
      parameters:
        - name: limit
          in: query
          description: How many items to return at one time (max 100)
          required: false
          schema:
            type: integer
            format: int32
      responses:
        '200':
          description: A paged array of pets
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Pets"
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
```
will generate code:
```
// Objects
type ErrorSchema struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type PetSchema struct {
	Id   int64   `json:"id"`
	Name string  `json:"name"`
	Tag  *string `json:"tag"`
}

type PetsSchema []PetSchema

// Parameters
type ListPetsParams struct {
	Limit *int32
}

// Response
type ListPetsResponse struct {
	Code        int
	Http200     PetsSchema
	HttpDefault *ErrorSchema
}

// Controller
type Controller interface {
	ListPets(params *ListPetsParams, req *http.Request, res http.ResponseWriter) ListPetsResponse
}

```

# Server boilerplate
If you want, library can generate echo server specific code. It will include only parameters validation, defaults and routes.

`oapi3gen -server echo spec.yaml`

```
func BuildRoutes(e *echo.Group, controller Controller) {
	e.GET("/pets", func(c echo.Context) error {
		parameters := &ListPetsParams{}
		if err := initParameters(c, parameters, nil); err != nil {
			return err
		}

		response := controller.ListPets(parameters, c.Request(), c.Response().Writer)
    ...
```
