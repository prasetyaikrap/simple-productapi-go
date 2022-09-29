package models

type (
	SuccessResponse struct {
		Code int
		Message string
		Data interface{}
	}

	ErrorResponse struct {
		Code int
		Message string
		Error string
	}
)