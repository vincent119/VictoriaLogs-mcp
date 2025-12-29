package server

import "fmt"

// Error definitions
var (
	// ErrInvalidRequest invalid request
	ErrInvalidRequest = fmt.Errorf("invalid request")

	// ErrMissingParameter missing required parameter
	ErrMissingParameter = fmt.Errorf("missing required parameter")

	// ErrQueryFailed query execution failed
	ErrQueryFailed = fmt.Errorf("query execution failed")

	// ErrServerNotReady server not ready
	ErrServerNotReady = fmt.Errorf("server not ready")
)

// ToolError tool error
type ToolError struct {
	Tool    string
	Message string
	Err     error
}

// Error implements error interface
func (e *ToolError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Tool, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Tool, e.Message)
}

// Unwrap implements errors.Unwrap
func (e *ToolError) Unwrap() error {
	return e.Err
}

// NewToolError creates a tool error
func NewToolError(tool, message string, err error) *ToolError {
	return &ToolError{
		Tool:    tool,
		Message: message,
		Err:     err,
	}
}
