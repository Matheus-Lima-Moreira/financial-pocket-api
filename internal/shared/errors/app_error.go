package errors

import "net/http"

type ErrorDetail struct {
	Code  string `json:"code"`
	Field string `json:"field"`
}

const (
	CodeRequired      = "REQUIRED"
	CodeInvalidEmail  = "INVALID_EMAIL"
	CodeMinLength     = "MIN_LENGTH"
	CodeMaxLength     = "MAX_LENGTH"
	CodeInvalidFormat = "INVALID_FORMAT"
	CodeInvalidJSON   = "INVALID_JSON"
	CodeMissingBody   = "MISSING_BODY"
	CodeNotFound      = "NOT_FOUND"
	CodeConflict      = "CONFLICT"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeBadRequest    = "BAD_REQUEST"
	CodeInternalError = "INTERNAL_ERROR"
	CodeInvalidToken  = "INVALID_TOKEN"
	CodeExpiredToken  = "EXPIRED_TOKEN"
	CodeMissingToken  = "MISSING_TOKEN"
)

type AppError struct {
	Code    int
	Message string
	Err     error
	ErrorCode string
	Field     string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.Code)
}

func (e *AppError) ToErrorDetail() ErrorDetail {
	code := e.ErrorCode
	if code == "" {
		code = CodeBadRequest
	}
	return ErrorDetail{
		Code:  code,
		Field: e.Field,
	}
}

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:      http.StatusNotFound,
		Message:   resource,
		ErrorCode: CodeNotFound,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:      http.StatusConflict,
		Message:   message,
		ErrorCode: CodeConflict,
	}
}

func NewUnauthorized(message string) *AppError {
	errorCode := CodeUnauthorized
	switch message {
	case "missing token":
		errorCode = CodeMissingToken
	case "invalid token":
		errorCode = CodeInvalidToken
	case "expired token":
		errorCode = CodeExpiredToken
	}

	return &AppError{
		Code:      http.StatusUnauthorized,
		Message:   message,
		ErrorCode: errorCode,
	}
}

func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:      http.StatusBadRequest,
		Message:   message,
		ErrorCode: CodeBadRequest,
	}
}

func NewValidationError(field string, code string) *AppError {
	return &AppError{
		Code:      http.StatusBadRequest,
		Message:   "validation error: " + field,
		ErrorCode: code,
		Field:     field,
	}
}
