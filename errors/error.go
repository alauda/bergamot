package errors

import "fmt"

// Code references to the kind of error that was returned
type Code string

// StatusMessage in order to make it simple
// to keep a list of the used codes and default messages
type StatusMessage struct {
	Message    string
	StatusCode int
}

// NewErrorStatusMessage build a message with a status code
func NewErrorStatusMessage(message string, status int) StatusMessage {
	return StatusMessage{Message: message, StatusCode: status}
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
	// ErrorCodeUnauthorized when user is not authorized
	ErrorCodeUnauthorized = "unauthorized"
	// ErrorCodeNotFound when is not found
	ErrorCodeNotFound = "not_found"
	// ErrorCodeElasticSearchError when using a elastic search and it returned an error
	ErrorCodeElasticSearchError = "elasticsearch_error"
	// ErrorCodeDatabaseError when using database commands and it returned an error
	ErrorCodeDatabaseError = "database_error"
)

var (
	// ErrorMessageList general default messages for all error codes
	ErrorMessageList = map[Code]StatusMessage{
		ErrorCodeResourceNotFound:      NewErrorStatusMessage("Resource not found", 404),
		ErrorCodeInvalidArgs:           NewErrorStatusMessage("Invalid parameters were passed.", 400),
		ErrorCodeBadRequest:            NewErrorStatusMessage("Required data not valid.", 400),
		ErrorCodeResourceAlreadyExists: NewErrorStatusMessage("The posted resource already existed.", 400),
		ErrorCodeUnknownIssue:          NewErrorStatusMessage("Unknown issue was caught and message was not specified.", 500),
		ErrorCodePermissionDenied:      NewErrorStatusMessage("Current user has no permission to perform the action.", 403),
		ErrorCodeResourceStateConflict: NewErrorStatusMessage("The posted resource already existed.", 409),
		ErrorCodeNotImplemented:        NewErrorStatusMessage("Method not implemented", 501),
		ErrorCodeElasticSearchError:    NewErrorStatusMessage("Elastic search error.", 500),
		ErrorCodeDatabaseError:         NewErrorStatusMessage("Database error.", 500),
	}
)

// AddError add error
func AddError(code Code, message string, status int) {
	ErrorMessageList[code] = NewErrorStatusMessage(message, status)
}

// AlaudaError common error struct used in the whole application
type AlaudaError struct {
	Source     string                `json:"source"`
	Message    string                `json:"message"`
	Code       Code                  `json:"code"`
	Fields     []map[string][]string `json:"fields,omitempty"`
	StatusCode int                   `json:"-"`
}

// New Constructor function for the error structure
func New(source string, code Code) *AlaudaError {
	var (
		message = "Error not described"
		status  = 999
	)
	if val, ok := ErrorMessageList[code]; ok {
		message = val.Message
		status = val.StatusCode
	}
	return &AlaudaError{
		Source:     source,
		Code:       code,
		Message:    message,
		StatusCode: status,
	}
}

// NewCommon return an error from a common error
func NewCommon(source string, err error) *AlaudaError {
	if err == nil {
		return nil
	}
	if alErr, ok := err.(*AlaudaError); ok {
		return alErr
	}
	return &AlaudaError{
		Source:     source,
		Code:       ErrorCodeUnknownIssue,
		Message:    err.Error(),
		StatusCode: ErrorMessageList[ErrorCodeUnknownIssue].StatusCode,
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

// SetCodeString sets a code as a string for custom codes
func (h *AlaudaError) SetCodeString(code string) *AlaudaError {
	return h.SetCode(Code(code))
}

// SetCode sets a code and returns itself for chaining calls
func (h *AlaudaError) SetCode(code Code) *AlaudaError {
	h.Code = code
	return h
}

// SetSource sets the source and returns itself for chaining calls
func (h *AlaudaError) SetSource(source string) *AlaudaError {
	h.Source = source
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

// GetError will create an error or return if already instatiated
func GetError(source string, err *AlaudaError, code Code) *AlaudaError {
	if err == nil {
		err = New(source, code)
	}
	return err
}

// TreatError treat common errors
func TreatError(source string, err error, code string) *AlaudaError {
	if alaudaErr, ok := err.(*AlaudaError); ok {
		return alaudaErr
	}
	return NewCommon(source, err).SetCodeString(code)
}
