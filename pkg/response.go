package pkg

import "github.com/go-playground/validator/v10"

type ApiResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type ErrorStruct struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Errors struct {
	Errors []ErrorStruct `json:"errors"`
}

func Null() interface{} {
	return nil
}

func BuildResponse[T any](statusCode int, message string, data T) ApiResponse[T] {
	return BuildResponse_(statusCode, message, data)
}

func BuildResponse_[T any](status int, message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}
}

func BuildErrorData(err error) Errors {

	var errors []ErrorStruct

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErr {
			errors = append(errors, ErrorStruct{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
	}

	return Errors{Errors: errors}
}
