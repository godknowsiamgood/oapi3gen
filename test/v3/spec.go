package v1

/**
    AUTOGENERATED. Please, do not edit.
**/

import (
	"net/http"
)

/* Components schemas */

type CoordinatesSchema struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type ImageSchema struct {
	AuthorName string  `json:"author_name,omitempty"`
	AuthorUrl  string  `json:"author_url,omitempty"`
	ImageUrl   string  `json:"image_url,omitempty"`
	Ratio      float64 `json:"ratio,omitempty"`
}

type PetAdviceSchema struct {
	Author   *PetAdviceAuthorSchema `json:"author,omitempty"`
	HtmlText string                 `json:"html_text"`
	Text     string                 `json:"text"`
}

type PetAdviceAuthorSchema struct {
	AvatarUrl    string `json:"avatar_url,omitempty"`
	InstagramUrl string `json:"instagram_url,omitempty"`
	Name         string `json:"name,omitempty"`
	Role         string `json:"role,omitempty"`
}

type PetAudioSchema struct {
	Duration      int64  `json:"duration,omitempty"`
	Transcription string `json:"transcription,omitempty"`
	Url           string `json:"url,omitempty"`
}

type PetBulletSchema struct {
	Color  string        `json:"color"`
	Emoji  string        `json:"emoji"`
	Images []ImageSchema `json:"images,omitempty"`
	Text   string        `json:"text"`
	Title  string        `json:"title"`
}

type PetButtonSchema struct {
	Coordinates *CoordinatesSchema `json:"coordinates,omitempty"`
	Title       string             `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Url         string             `json:"url,omitempty"`
}

type WalkPetSchema struct {
	Address          string              `json:"address,omitempty"`
	Coordinates      *CoordinatesSchema  `json:"coordinates,omitempty"`
	Description      string              `json:"description,omitempty"`
	DescriptionShort string              `json:"description_short,omitempty"`
	Id               int64               `json:"id,omitempty"`
	Image            *ImageSchema        `json:"image,omitempty"`
	IsKey            bool                `json:"is_key"`
	Path             []CoordinatesSchema `json:"path"`
	Title            string              `json:"title,omitempty"`
	Type             string              `json:"type"`
}

type WalkPreviewSchema struct {
	Description string       `json:"description,omitempty"`
	Iata        string       `json:"iata,omitempty"`
	Id          int64        `json:"id,omitempty"`
	Image       *ImageSchema `json:"image,omitempty"`
	IsSoon      bool         `json:"is_soon,omitempty"`
	Title       string       `json:"title,omitempty"`
}

/* Components responses */

type ErrorResponse struct {
	Message string `json:"message"`
}

/* Parameters */

type GetPetByIdParams struct {
	Id     int64
	Locale string
}

type GetPetAudioByIdParams struct {
	Id     int64
	Locale string
}

type SavePetBulletParams struct {
}

type SaveImageParams struct {
}

type GetStyleParams struct {
	Theme string
}

type GetWalkByIdParams struct {
	Id     int64
	Locale string
}

type GetWalksPreviewParams struct {
	Iata   string
	Locale string
}

/* Requests bodies */

type SavePetBulletBody struct {
	PetId  *int64  `json:"Pet_id"`
	Color  string  `json:"color"`
	Emoji  string  `json:"emoji"`
	Images []int64 `json:"images"`
}

type SaveImageBody struct {
}

/* Response objects */

type GetPetByIdHttp200Response struct {
	Address          string            `json:"address,omitempty"`
	Advice           *PetAdviceSchema  `json:"advice,omitempty"`
	Audio            *PetAudioSchema   `json:"audio,omitempty"`
	Bullets          []PetBulletSchema `json:"bullets,omitempty"`
	Buttons          []PetButtonSchema `json:"buttons"`
	Description      string            `json:"description,omitempty"`
	DescriptionShort string            `json:"description_short,omitempty"`
	Id               int64             `json:"id,omitempty"`
	Images           []ImageSchema     `json:"images,omitempty"`
	Title            string            `json:"title,omitempty"`
}

type SaveImageHttp200Response struct {
	ImageId int64 `json:"image_id"`
}

type GetWalkByIdHttp200Response struct {
	Pets        []WalkPetSchema `json:"Pets,omitempty"`
	Description string          `json:"description,omitempty"`
	Duration    int64           `json:"duration,omitempty"`
	Id          int64           `json:"id,omitempty"`
	Image       *ImageSchema    `json:"image,omitempty"`
	Length      int64           `json:"length,omitempty"`
	MapGeoJson  string          `json:"map_geo_json,omitempty"`
	MapImageUrl string          `json:"map_image_url,omitempty"`
	MaxZoom     float64         `json:"max_zoom,omitempty"`
	MinZoom     float64         `json:"min_zoom,omitempty"`
	StartZoom   float64         `json:"start_zoom,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Title       string          `json:"title,omitempty"`
}

type GetWalksPreviewHttp200Response struct {
	Walks []WalkPreviewSchema `json:"walks,omitempty"`
}

/* Responses */

type GetPetByIdResponse struct {
	Code    int
	Http200 *GetPetByIdHttp200Response
}

type SaveImageResponse struct {
	Code    int
	Http200 *SaveImageHttp200Response
}

type GetStyleResponse struct {
	Code    int
	Http200 map[string]interface{}
}

type GetWalkByIdResponse struct {
	Code    int
	Http200 *GetWalkByIdHttp200Response
}

type GetWalksPreviewResponse struct {
	Code    int
	Http200 *GetWalksPreviewHttp200Response
}

type Controller interface {
	GetPetById(params *GetPetByIdParams, req *http.Request, res http.ResponseWriter) (GetPetByIdResponse, error)
	GetPetAudioById(params *GetPetAudioByIdParams, req *http.Request, res http.ResponseWriter) (int, error)
	SavePetBullet(params *SavePetBulletParams, body *SavePetBulletBody, req *http.Request, res http.ResponseWriter) (int, error)
	SaveImage(params *SaveImageParams, body *SaveImageBody, req *http.Request, res http.ResponseWriter) (SaveImageResponse, error)
	GetStyle(params *GetStyleParams, req *http.Request, res http.ResponseWriter) (GetStyleResponse, error)
	GetWalkById(params *GetWalkByIdParams, req *http.Request, res http.ResponseWriter) (GetWalkByIdResponse, error)
	GetWalksPreview(params *GetWalksPreviewParams, req *http.Request, res http.ResponseWriter) (GetWalksPreviewResponse, error)

	Error(err error) ErrorResponse
}
