package dto

type BadRequestResponse struct {
	Errors string `json:"errors"`
}

type UnauthorizedResponse struct {
	Errors string `json:"errors"`
}

type InternalServerErrorResponse struct {
	Errors string `json:"errors"`
}
