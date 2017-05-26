package mekansm

import "fmt"

// ErrorCode references to the kind of error that was returned
type ErrorCode string

// ErrorStatusMessage in order to make it simple
// to keep a list of the used codes and default messages
type ErrorStatusMessage struct {
	Message    string
	StatusCode int
}

// NewErrorStatusMessage build a message with a status code
func NewErrorStatusMessage(message string, status int) ErrorStatusMessage {
	return ErrorStatusMessage{Message: message, StatusCode: status}
}

const (
	// ErrorCodeResourceNotFound error when the requested resource was not found
	ErrorCodeResourceNotFound = "resource_not_exist"
	// ErrorCodeInvalidArgs error when the request data is incorrect or incomplete
	ErrorCodeInvalidArgs = "invalid_args"
	// ErrorCodeBadRequest error when the needed input is not provided
	ErrorCodeBadRequest = "bad_request"
	//ErrorCodeResourceAlreadyExists when a resource already exists
	ErrorCodeResourceAlreadyExists = "resource_already_exist"
	// ErrorCodeUnknownIssue when the issue is unknown
	ErrorCodeUnknownIssue = "unknown_issue"
	// ErrorCodePermissionDenied when user dont have necessary permissions
	ErrorCodePermissionDenied = "permission_denied"
	// ErrorCodeResourceStateConflict when the resource is in another state and generates a conflict
	ErrorCodeResourceStateConflict = "resource_state_conflict"
	// ErrorCodeNotImplemented when some method is not implemented
	ErrorCodeNotImplemented = "not_implemented"
)

var (
	// ErrorMessageList general default messages for all error codes
	ErrorMessageList = map[ErrorCode]ErrorStatusMessage{
		ErrorCodeResourceNotFound:      NewErrorStatusMessage("Resource not found", 404),
		ErrorCodeInvalidArgs:           NewErrorStatusMessage("Invalid parameters were passed.", 400),
		ErrorCodeBadRequest:            NewErrorStatusMessage("Required data not valid.", 400),
		ErrorCodeResourceAlreadyExists: NewErrorStatusMessage("The posted resource already existed.", 400),
		ErrorCodeUnknownIssue:          NewErrorStatusMessage("Unknown issue was caught and message was not specified.", 500),
		ErrorCodePermissionDenied:      NewErrorStatusMessage("Current user has no permission to perform the action.", 403),
		ErrorCodeResourceStateConflict: NewErrorStatusMessage("The posted resource already existed.", 409),
		ErrorCodeNotImplemented:        NewErrorStatusMessage("Method not implemented", 501),
	}
)

// AlaudaError common error struct used in the whole application
type AlaudaError struct {
	Source     string                `json:"source"`
	Message    string                `json:"message"`
	Code       ErrorCode             `json:"code"`
	Fields     []map[string][]string `json:"fields,omitempty"`
	StatusCode int                   `json:"-"`
}

// NewAlaudaError Constructor function for the error structure
func NewAlaudaError(source string, code ErrorCode) *AlaudaError {
	return &AlaudaError{
		Source:     source,
		Code:       code,
		Message:    ErrorMessageList[code].Message,
		StatusCode: ErrorMessageList[code].StatusCode,
	}
}

// Error satisfies the error interface
func (h *AlaudaError) Error() string {
	return fmt.Sprintf("%s: %s", h.Code, h.Message)
}

// SetMessage sets a message using the given format and parameters
// returns itself
func (h *AlaudaError) SetMessage(format string, args ...interface{}) *AlaudaError {
	h.Message = fmt.Sprintf(format, args...)
	return h
}

// AddFieldError adds a field error and return itself
func (h *AlaudaError) AddFieldError(field string, message ...string) *AlaudaError {
	h.initializeFields()
	if value, ok := h.Fields[0][field]; ok {
		value = append(value, message...)
		return h
	}
	h.Fields[0][field] = message
	return h
}

// AddFieldsFromError will add the field errors of the given error if any
// will also add in the current error using a formated string as base
func (h *AlaudaError) AddFieldsFromError(format string, index int, err *AlaudaError) *AlaudaError {
	if err == nil || err.Fields == nil || len(err.Fields) == 0 || len(err.Fields[0]) == 0 {
		return h
	}
	for k, v := range err.Fields[0] {
		h.AddFieldError(fmt.Sprintf(format, index, k), v...)
	}
	return h
}

func (h *AlaudaError) initializeFields() {
	if h.Fields == nil {
		h.Fields = []map[string][]string{
			map[string][]string{},
		}
	}
}
