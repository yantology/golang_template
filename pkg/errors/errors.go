package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Client errors (4xx)
	ErrorCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrorCodeConflict            ErrorCode = "CONFLICT"
	ErrorCodeValidation          ErrorCode = "VALIDATION_ERROR"
	ErrorCodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"

	// Server errors (5xx)
	ErrorCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadGateway     ErrorCode = "BAD_GATEWAY"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeGatewayTimeout  ErrorCode = "GATEWAY_TIMEOUT"

	// Business logic errors
	ErrorCodeBusinessLogic   ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrorCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// AppError represents an application error with context
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	StatusCode int                    `json:"-"`
	Cause      error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// MarshalJSON customizes JSON marshaling for AppError
func (e *AppError) MarshalJSON() ([]byte, error) {
	type alias AppError
	return json.Marshal(&struct {
		*alias
		Error string `json:"error,omitempty"`
	}{
		alias: (*alias)(e),
		Error: e.Error(),
	})
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// Newf creates a new AppError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Cause:      err,
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		Cause:      err,
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithField adds a field to the error context
func (e *AppError) WithField(key string, value interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[key] = value
	return e
}

// WithFields adds multiple fields to the error context
func (e *AppError) WithFields(fields map[string]interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// WithStatusCode sets a custom HTTP status code
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	e.StatusCode = statusCode
	return e
}

// GetStatusCode returns the HTTP status code for the error
func (e *AppError) GetStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	return getStatusCodeForErrorCode(e.Code)
}

// IsType checks if the error is of a specific type
func (e *AppError) IsType(code ErrorCode) bool {
	return e.Code == code
}

// Predefined error constructors
func NewBadRequestError(message string) *AppError {
	return New(ErrorCodeBadRequest, message)
}

func NewUnauthorizedError(message string) *AppError {
	return New(ErrorCodeUnauthorized, message)
}

func NewForbiddenError(message string) *AppError {
	return New(ErrorCodeForbidden, message)
}

func NewNotFoundError(message string) *AppError {
	return New(ErrorCodeNotFound, message)
}

func NewConflictError(message string) *AppError {
	return New(ErrorCodeConflict, message)
}

func NewValidationError(message string) *AppError {
	return New(ErrorCodeValidation, message)
}

func NewInternalServerError(message string) *AppError {
	return New(ErrorCodeInternalServer, message)
}

func NewBusinessLogicError(message string) *AppError {
	return New(ErrorCodeBusinessLogic, message)
}

func NewDatabaseError(err error) *AppError {
	return Wrap(err, ErrorCodeDatabaseError, "Database operation failed")
}

func NewExternalServiceError(service string, err error) *AppError {
	return Wrap(err, ErrorCodeExternalService, fmt.Sprintf("External service %s failed", service))
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Code
	}
	return ErrorCodeInternalServer
}

// GetStatusCode extracts the HTTP status code from an error
func GetStatusCode(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.GetStatusCode()
	}
	return http.StatusInternalServerError
}

// getStatusCodeForErrorCode maps error codes to HTTP status codes
func getStatusCodeForErrorCode(code ErrorCode) int {
	switch code {
	case ErrorCodeBadRequest:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeForbidden:
		return http.StatusForbidden
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorCodeValidation:
		return http.StatusBadRequest
	case ErrorCodeTooManyRequests:
		return http.StatusTooManyRequests
	case ErrorCodeUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrorCodeBadGateway:
		return http.StatusBadGateway
	case ErrorCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrorCodeGatewayTimeout:
		return http.StatusGatewayTimeout
	case ErrorCodeBusinessLogic:
		return http.StatusBadRequest
	case ErrorCodeDatabaseError:
		return http.StatusInternalServerError
	case ErrorCodeExternalService:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}